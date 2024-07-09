image=$(docker images --filter=reference='rarime-points-svc:test' -q)

if [[ -n $image ]]; then
    echo "[+] Remove image $image ~ rarime-points-svc:test"
    docker rmi -f $image > /dev/null
    echo "[+] Image successfully removed."
fi

echo "[+] Building image for tests..."
docker build --tag rarime-points-svc:test --file ./Dockerfile.testing .
if [[ $? -ne 0 ]]; then
    echo "[-] Error while building!"
    exit $?
fi
echo "[+] Image for tests successfully built: $(docker images --filter=reference='rarime-points-svc:test' -q)"
echo "[+] Run docker compose test environment"
docker compose -f docker-compose.testing.yaml up -d

sleep 5

echo "[+] Run tests"
KV_VIPER_FILE=./tests/config-testing-0-external.yaml go test requests_test.go

docker compose -f docker-compose.testing.yaml down -v
