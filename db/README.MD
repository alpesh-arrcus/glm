

TableName: customers
Schema: customerName (text), customerSecret (text), id (int), status(text)
Status Values: "Active", "Inactive"


TableName: {customerName}_devices
Schema: fp (text), lastHB (text), status (text)
fp - device fingerprint - has to be unique
lastHB - string representation of a time at which last heart beat was gotten from this device
status - "Active", "RMA"


//Assumptions - licenses for same feature will have same usagePeriod
TableName:  {customerName}_purchases
Schema: featureName(text), licCount(int), usagePeriod (int), purchaseTime(text)
licCount - count remaining on a purchased license [Constraint >= 0]
usagePeriod - seconds of usage for each of the license


TableName: {customerName}_licenseAllocs
Schema: featureName(text), deviceFp(text), status(text), periodLeft(int), lastUse(text)
deviceFp - is "" when the license is NOT allocated to any device
status - "InUse", "Available"
