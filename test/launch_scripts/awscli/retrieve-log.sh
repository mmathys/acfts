# retrieves logs from a validator

# get log file location
ssh -i ~/.ssh/makesxi-frankfurt ec2-user@$1 "docker logs acfts_server_1"