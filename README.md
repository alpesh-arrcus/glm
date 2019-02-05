# glm
Global License Mgr


Steps to run:
```
$ wget https://raw.githubusercontent.com/toravir/glm/master/docker/Dockerfile

$ docker build -t try1 .
...
... #Takes about two minutes

$ docker run -it try1 bash
root@c9c25bd2bd56:/go# cd bin          
root@c9c25bd2bd56:/go/bin# ./sanity.sh 
```
