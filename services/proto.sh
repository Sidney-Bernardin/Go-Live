function generate_for {
    for svc in $@
    do
        protoc \
            --go_out="$svc"/proto \
            --go-grpc_out="$svc"/proto \
            "$1"/proto/"$1".proto

        echo "Generated $1""-service's gPRC files for $svc""-service."
    done 
}

for arg in $@
do
    case $arg in
        'users' )
            generate_for $arg 'rooms'
            ;;
    esac
done
