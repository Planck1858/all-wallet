up:
	docker-compose up -d --build

copy-to-server:
	echo "copy to server ${SERVER_IP}"
	scp -r [!.]* root@${SERVER_IP}:/root/go/src/github.com/Planck1858/all-wallet
