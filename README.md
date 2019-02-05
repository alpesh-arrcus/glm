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
1. Add a customer's purchase of license
{"statusCode":200,"statusString":"Sucessful"}
2. Check the db                              
c1|c1123|0|Active                            
feat1|10|100|2019-02-05T16:33:05Z            
3. Register a device                         
{"statusCode":201,"statusString":"Successful","curTime":"2019-02-05T16:33:05Z"}
4. Allocate a license                                                          
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:33:05Z"}
5. Verify Allocated license is removed from purchases and moved to allocs      
feat1|9|100|2019-02-05T16:33:05Z                                               
feat1|device1|InUse|100|2019-02-05T16:33:05Z                                   
6. Sleep for 3s                                                                
7. Punch HB                                                                    
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:33:08Z","expiredLicenses":[]}
8. Verify periodLeft go down                                                                        
feat1|device1|InUse|97|2019-02-05T16:33:05Z                                                         
9. Sleep for 30s                                                                                    
10. Punch HB                                                                                        
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:33:38Z","expiredLicenses":[]}
11. Verify periodLeft go down                                                                       
feat1|device1|InUse|67|2019-02-05T16:33:05Z                                                         
12. Sleep for 75s                                                                                   
13. Punch HB                                                                                        
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:34:53Z","expiredLicenses":[]}
14. Verify periodLeft go down to zero and another license allocated
feat1|device1|Expired|0|2019-02-05T16:33:05Z
feat1|device1|InUse|100|2019-02-05T16:34:53Z
feat1|8|100|2019-02-05T16:33:05Z
15. Sleep for 100s
16. Punch HB - saying don't allocate expiring license
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:36:33Z","expiredLicenses":[]}
17. Punch HB - saying don't allocate expiring license
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:36:33Z","expiredLicenses":[]}
18. Verify periodLeft go down to zero and NO new license allocated
feat1|device1|Expired|0|2019-02-05T16:33:05Z
feat1|device1|InUse|0|2019-02-05T16:34:53Z
feat1|8|100|2019-02-05T16:33:05Z
19. Allocate a license
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:36:33Z"}
20. Verify Allocated license is removed from purchases and moved to allocs
feat1|7|100|2019-02-05T16:33:05Z
feat1|device1|Expired|0|2019-02-05T16:33:05Z
feat1|device1|InUse|0|2019-02-05T16:34:53Z
feat1|device1|InUse|100|2019-02-05T16:36:33Z
21. Punch HB
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:36:33Z","expiredLicenses":[]}
22. Sleep for 30s
23. Punch HB
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:37:03Z","expiredLicenses":[]}
24. Sleep for 30s
25. Release the license
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:37:33Z"}
26. Verify license Status is not InUse
feat1|device1|Available|70|2019-02-05T16:37:03Z
feat1|device1|Available|70|2019-02-05T16:37:03Z
feat1|device1|Available|70|2019-02-05T16:37:03Z
27. Sleep for 50s
28. Punch HB verify expiredLicenses stays empty
{"statusCode":200,"statusString":"Successful","curTime":"2019-02-05T16:43:12Z","expiredLicenses":[]}
29. Verify license Status is not InUse and periodLeft stays same
feat1|device1|Available|70|2019-02-05T16:41:52Z
feat1|device1|Available|70|2019-02-05T16:41:52Z
feat1|device1|Available|70|2019-02-05T16:41:52Z
30. Done
```
