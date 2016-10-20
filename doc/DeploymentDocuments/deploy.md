auto-deploy

1. checkEnvironment
2. deploy

1. checkEnvironment
   * check git, docker, kubernetes, golang
   * check images (mysql, redis-master, redis-slave)
   * git connection validate 
   
2. deploy
   * create ns, quota and limit
   * git init
   * git pull frontend
   * git pull backend
   
   * build docker image
   * docker rmi
   * docker push
   
   * deploy mysql
   * deploy redis-master, redis-slave
   * init mysql
   * deploy yce
   
   
update
1. update preparation
   * git pull frontend
   * git pull backend
   * docker rmi
   * docker push
   
2. update manually 