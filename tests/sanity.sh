#!/bin/bash

DBPATH=/data/glm.db

rm -f $DBPATH
pkill glm

./glm -config /tmp/example_cfg.ini &

sleep 1

echo "1. Add a customer's purchase of license"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/addPurchase -d '{"customerName":"c1", "featureName" : "feat1", "licenseCount":10, "usagePeriod" : 100}'


echo "2. Check the db"
echo "select * from customers;" | sqlite3 $DBPATH
echo "select * from c1_purchases;" | sqlite3 $DBPATH

echo "3. Register a device"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/deviceInit -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "bootTime":"sometime"}'

echo "4. Allocate a license"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/licenseAlloc -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "featureName" : "feat1" }'

echo "5. Verify Allocated license is removed from purchases and moved to allocs"
echo "select * from c1_purchases;" | sqlite3 $DBPATH
echo "select * from c1_licenseAllocs;" | sqlite3 $DBPATH

echo "6. Sleep for 3s"
sleep 3

echo "7. Punch HB"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/deviceHB -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "curTime" :"soemtime", "autoReAllocExpiring": true}'

echo "8. Verify periodLeft go down"
echo "select * from c1_licenseAllocs;" | sqlite3 $DBPATH

echo "9. Sleep for 30s"
sleep 30

echo "10. Punch HB"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/deviceHB -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "curTime" :"soemtime", "autoReAllocExpiring": true}'

echo "11. Verify periodLeft go down"
echo "select * from c1_licenseAllocs;" | sqlite3 $DBPATH

echo "12. Sleep for 75s"
sleep 75

echo "13. Punch HB"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/deviceHB -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "curTime" :"soemtime", "autoReAllocExpiring": true}'

echo "14. Verify periodLeft go down to zero and another license allocated"
echo "select * from c1_licenseAllocs;" | sqlite3 $DBPATH
echo "select * from c1_purchases;" | sqlite3 $DBPATH

echo "15. Sleep for 100s"
sleep 100

echo "16. Punch HB - saying don't allocate expiring license"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/deviceHB -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "curTime" :"soemtime", "autoReAllocExpiring": false}'

echo "17. Punch HB - saying don't allocate expiring license"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/deviceHB -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "curTime" :"soemtime", "autoReAllocExpiring": false}'

echo "18. Verify periodLeft go down to zero and NO new license allocated"
echo "select * from c1_licenseAllocs;" | sqlite3 $DBPATH
echo "select * from c1_purchases;" | sqlite3 $DBPATH

echo "19. Allocate a license"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/licenseAlloc -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "featureName" : "feat1" }'

echo "20. Verify Allocated license is removed from purchases and moved to allocs"
echo "select * from c1_purchases;" | sqlite3 $DBPATH
echo "select * from c1_licenseAllocs;" | sqlite3 $DBPATH

echo "21. Punch HB"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/deviceHB -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "curTime" :"soemtime", "autoReAllocExpiring": true}'

echo "22. Sleep for 30s"
sleep 30

echo "23. Punch HB"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/deviceHB -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "curTime" :"soemtime", "autoReAllocExpiring": true}'

echo "24. Sleep for 30s"
sleep 30

echo "25. Release the license"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/licenseFree -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "featureName" : "feat1", "curTime":"isignored" }'

echo "26. Verify license Status is not InUse"
echo "select * from c1_licenseAllocs;" | sqlite3 $DBPATH

echo "27. Sleep for 50s"
sleep 50

echo "28. Punch HB verify expiredLicenses stays empty"
curl -k -H 'Content-Type: application/json'  https://localhost:11223/c1/deviceHB -d '{"customerName":"c1", "customerSecret":"c1123", "fingerPrint":"device1", "curTime" :"soemtime", "autoReAllocExpiring": true}'

echo "29. Verify license Status is not InUse and periodLeft stays same"
echo "select * from c1_licenseAllocs;" | sqlite3 $DBPATH

echo "30. Done"
