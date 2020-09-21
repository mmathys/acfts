#!/bin/bash
amazon-linux-extras install docker -y
service docker start
usermod -a -G docker ec2-user
chkconfig docker on
yum install -y git htop
curl -L https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
echo -e "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAviUCscfHjgVfPrdtlt8bsjqxZ/U0SwIco28hKevQdh4mVmBy\n3EbMDQnkNStyC8Uvqdvm8jfyMVAM7K1RUkXXWusz9fIVhzirNoOgx32VymD44bGN\nV/RM4KWEbnERmzGD6djjmk8K9pmSmBojD2DEhdREnfeko2jLhdHIPV5qApOKkBWs\nmuDV1bEhqtVulCTAhukQZYDYOQkY4vuWzKXEOk2HAul5A0/ME3pKkw4xr5bdh8kd\nMyZMtYZjsGxQWsNOvfTcSqAyqji3Pa6WEixHuxjhaQ4LGMulVorM+l8R7umxJE32\no0ISL8QbPVwws7ELrs0GtntUKXGoa/sICtcaBwIDAQABAoIBAQCGHTK+VGwjdxSL\n9vL8NFejqZRR//inY1Dx0cNxfxulBHUD/blmNyVbmPt2qTMirC2HQEyv0dgpAg5r\nrX+raQmR3k3jDqIyozzYEZWIv8K0YRuKjSJlxDGpTfWPtH63mQAzI8C9qTzRXNay\nI9VQ1m8m+eaACRj+DT9oHuR+rKZXM6NnMTEZlyzlvW6UfvPSHgAXABh4Ith6GgBB\nZq0sMOH+Y1lP+6bHotkYVeoMlk5mWlqZJC60PUHh9l7cISRygemVhR2sZPgcVeut\nWtBWHX61tEQoRCm78k4hToDPWYRYu21I5g8uZoOfQaLSuKMy6+rZ+tj5MfSIVbmV\nWOyVb7oRAoGBAPB5OGN45jDb36J8jfL9zjnAMdMLn0Gmk24yKMy2MAjvQkpn61nY\nQ3/GLidIR8hh6f9bYc8luaRrvX6ZV+ROTQk6P6EntR2wXQz0QYMQ/fpumMlK1JD2\nQo0h+CzF5NKQ9t7C4xFP1elsMH50eVSu/u8iwRN7JXFohBFmqo1gRqd/AoGBAMpr\n55RWEEu6uHqzabN7s1mS8ypGyLOz9QZu4Zsz4jeM9lwSPWbXTXHOHT6QiXExDM+I\njyh6vHpNzdSRMKFOJdL9N1TnCVbn3ru1T8d8sXIgVBEda3KXY8t6C4oDt9ELzBOz\nckiBc5oz5kYsBCCfIDX04SyDLS6YwMlzV9WM7ZF5AoGAGZI6XdGIYN+wqz2ldpZo\nrJl0cuXEfJ2zBb/NuN8SwpB2syTT0Qz/OIIhfEniYDuRP3c/MmDnmcVfNCwoPX3a\nEUYhaAip6MtNSchbgw12hH/oPVd4vSWCtDFNi61LhNync86ziK69p4G6sme+UzP8\npf9DGsg4lhOwH42US05WlQ8CgYEAtq3FFDjLghyj6pBKH2HOzVLyQsAY1xzcT473\n6gJb8Lg1iYgog8xgpwlVmSYF8/JJVm4ujGCgAlx/Goxhb+G//vbo3aWX5ydg+tQ5\niJX+6jlExSFg28zGkla1XTvXw0A7T7hNtiOluY+HDOMiYnJZNnimUm6b8waGTgfb\nR+1kZTkCgYAh971eMcyo9juN44WqOmT+7qv14ba1prXxvMShgdYGnygr0omfRGRO\nofOo5c5JLyKuBAF9fckGjHS4//QiqfJwOO4kJKS02WNgqF+gRYHpRVujsCscIeXC\nRdTM361NSLrwS6b5sf9wwKGlboFj42T9wTxQZhQhGSRJpHZHtkyawA==\n-----END RSA PRIVATE KEY-----" >> /home/ec2-user/.ssh/id_rsa
echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+JQKxx8eOBV8+t22W3xuyOrFn9TRLAhyjbyEp69B2HiZWYHLcRswNCeQ1K3ILxS+p2+byN/IxUAzsrVFSRdda6zP18hWHOKs2g6DHfZXKYPjhsY1X9EzgpYRucRGbMYPp2OOaTwr2mZKYGiMPYMSF1ESd96SjaMuF0cg9XmoCk4qQFaya4NXVsSGq1W6UJMCG6RBlgNg5CRji+5bMpcQ6TYcC6XkDT8wTekqTDjGvlt2HyR0zJky1hmOwbFBaw0699NxKoDKqOLc9rpYSLEe7GOFpDgsYy6VWisz6XxHu6bEkTfajQhIvxBs9XDCzsQuuzQa2e1Qpcahr+wgK1xoH max@mmathys" >> /home/ec2-user/.ssh/id_rsa.pub

cp /home/ec2-user/.ssh/id_rsa /root/.ssh/
cp /home/ec2-user/.ssh/id_rsa.pub /root/.ssh/
chmod 600 /root/.ssh/id_rsa
chmod 644 /root/.ssh/id_rsa.pub

chown ec2-user /home/ec2-user/.ssh/id_rsa
chown ec2-user /home/ec2-user/.ssh/id_rsa.pub
chmod 600 /home/ec2-user/.ssh/id_rsa
chmod 644 /home/ec2-user/.ssh/id_rsa.pub

echo "github.com,140.82.118.3 ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==" >> /home/ec2-user/.ssh/known_hosts

cp /home/ec2-user/.ssh/known_hosts /root/.ssh/known_hosts
chmod 644 /root/.ssh/known_hosts

chown ec2-user /home/ec2-user/.ssh/known_hosts
chmod 644 /home/ec2-user/.ssh/known_hosts

git clone git@github.com:mmathys/acfts.git /home/ec2-user/acfts
chown -R ec2-user /home/ec2-user/acfts
cd /home/ec2-user/acfts
