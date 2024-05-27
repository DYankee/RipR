album info
    get metadata
    set album info
    build label file
    import label file


audio manipulation
    trim album sides
    clean pops
    
    split tracks on labels
    export flac files


The Program:

command line tool to split and name records that have been recorded on audacity

rriper -a "artist" -r "release" -f "audacity project file"

1. open audacity project
    
2. search music brainz for releases
3. return list of vinyl records
4. get user choice
5. get song info from music brainz
6. calculate song lengths based on recorded side length
7. split songs in audacity acording to calculated length
8. export each song as a flec file


Current todo:

add audacity commands
    AddLabel:
