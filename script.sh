function kill_process() {
    if [[ -f ".rose.pid" ]]; then
        kill $(cat .rose.pid)
    fi
}

kill_process
