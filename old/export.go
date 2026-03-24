package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
)

// calcLengthMod computes the proportional time adjustment factor
// to align MusicBrainz track durations with the actual audio clip
// length in Audacity.
func (sd *sideData) calcLengthMod() error {
	if sd.songExportData == nil {
		return fmt.Errorf("no song export data for side")
	}

	sd.sideLength = 0
	for _, s := range sd.songExportData {
		sd.sideLength += s.songLength // milliseconds
	}
	sd.sideLength /= 1000 // convert to seconds

	clipLength, err := sd.clipInfo.GetClipLength()
	if err != nil {
		return err
	}

	log.Printf(
		"MB side length: %.3fs, Audacity clip length: %.3fs",
		sd.sideLength,
		clipLength,
	)

	// Ratio: how much longer/shorter the clip is vs MusicBrainz data.
	// We scale each track by this ratio so the sum matches the clip.
	sd.lengthMod = clipLength / sd.sideLength

	log.Printf("Length mod (scale factor): %f", sd.lengthMod)
	return nil
}

func (m *model) buildExportData() {
	// BUG FIX: distribute tracks across sides based on medium index,
	// not dumping everything into sideData[0].
	for i, med := range m.releaseData.Mediums {
		if i >= len(m.sideData) {
			log.Printf(
				"warning: more mediums (%d) than clips (%d), skipping",
				len(m.releaseData.Mediums),
				len(m.sideData),
			)
			break
		}
		for _, t := range med.Tracks {
			m.sideData[i].songExportData = append(
				m.sideData[i].songExportData,
				songData{
					songName:     t.Recording.Title,
					songLength:   float64(t.Recording.Length),
					songPosition: t.Position,
				},
			)
		}
		if err := m.sideData[i].calcLengthMod(); err != nil {
			log.Printf("calcLengthMod error for side %d: %v", i, err)
		}
	}
}

func (m *model) exportSongs() error {
	outDir := m.inputs[2].Value()
	if outDir == "" {
		outDir = "output"
	}

	if err := os.MkdirAll(outDir, 0700); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	absDir, err := filepath.Abs(outDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output dir: %w", err)
	}

	for sideIdx, sd := range m.sideData {
		if len(sd.songExportData) == 0 {
			continue
		}

		offset := float64(sd.clipInfo.Start)
		log.Printf("--- Exporting side %d (offset start: %.3f) ---",
			sideIdx, offset)

		for _, s := range sd.songExportData {
			// BUG FIX: the original code modified the range-copy of s,
			// then recalculated differently in the second loop.
			// Now we do it once, correctly.
			// Convert ms → seconds, then scale by the length mod.
			scaledLength := (s.songLength / 1000) * sd.lengthMod

			// Add a small buffer (0.5s) to avoid cutting off endings.
			regionEnd := offset + scaledLength + 0.5

			log.Printf(
				"Track %d %q: offset=%.3f length=%.3f end=%.3f",
				s.songPosition, s.songName,
				offset, scaledLength, regionEnd,
			)

			if err := m.audacity.Tracks.SelectRegion(
				offset, math.Min(regionEnd, offset+scaledLength+1),
			); err != nil {
				return fmt.Errorf("SelectRegion failed: %w", err)
			}

			filename := fmt.Sprintf(
				"%02d_%s.flac",
				s.songPosition,
				sanitizeFilename(s.songName),
			)

			if err := m.audacity.IO.ExportAudio(
				absDir, filename,
			); err != nil {
				log.Printf("export error for %s: %v", filename, err)
			}

			offset += scaledLength
		}
	}

	return nil
}

// sanitizeFilename removes or replaces characters that are
// problematic in file paths.
func sanitizeFilename(name string) string {
	var result []rune
	for _, r := range name {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			result = append(result, '_')
		default:
			result = append(result, r)
		}
	}
	return string(result)
}
