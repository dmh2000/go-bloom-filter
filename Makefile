redis:
	sudo docker run -d --name redis-stack -p 6379:6379 redis/redis-stack:latest

docker-kill:
	-sudo docker ps -a -q | xargs sudo docker kill

docker-rm: docker-kill
	-sudo docker  container list -aq | xargs sudo docker container rm

docker-clean:
	-sudo docker ps -a -q | xargs sudo docker kill
	-sudo docker  container list -aq | xargs sudo docker container rm
	-sudo docker  image list -q | xargs sudo docker image rm