# ansible-player

This is a simple daemon to run Ansible playbook from HTTP calls

Define you job in a file like this
```
site:
  git: git@github.com:wayt/ansible-player.git
  command: ansible-playbook -i hosts site.yml
```

ansible-player will use the container ssh key to clone / access to the servers.

* `docker run -d --name ansible-player --restart=always -p 8080:8080 -v /your/.ssh/directory:/root/.ssh -v /your/jobs/file:/root/jobs:ro -v /your/access/file:/root/access:ro -v /your/logs/directory:/root/logs -e LOG_DIR=/root/logs -e AUTH_FILE=/root/access -e JOB_FILE=/root/jobs maxwayt/ansible-player`

To run a playbook with curl:

* `curl -X POST -u username:password your_host:8080/job -d name=site`

The access file is similar to htpasswd, with sha1. Format:
```
username1:sha1(password)
username2:sha1(password)
```
