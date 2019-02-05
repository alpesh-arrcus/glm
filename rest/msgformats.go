package rest

type deviceInitReq struct {
	Fingerprint    string `json:fingerPrint`
	BootupTime     string `json:bootTime`
	CustomerName   string `json:customerName`
	CustomerSecret string `json:customerSecret`
}

type deviceInitResp struct {
	Status       int    `json:statusCode`
	StatusString string `json:statusString`
	CurTime      string `json:curTime`
}

type licenseAllocReq struct {
	FeatureName    string `json:featureName`
	Fingerprint    string `json:fingerPrint`
	CustomerName   string `json:customerName`
	CustomerSecret string `json:customerSecret`
}

type licenseAllocResp struct {
	Status       int    `json:statusCode`
	StatusString string `json:statusString`
	CurTime      string `json:curTime`
}

type licenseFreeReq struct {
	Fingerprint    string `json:fingerPrint`
	curTime        string `json:curTime`
	CustomerName   string `json:customerName`
	CustomerSecret string `json:customerSecret`
}

type licenseFreeResp struct {
	Status       int    `json:statusCode`
	StatusString string `json:statusString`
	curTime      string `json:curTime`
}

type deviceHBPunchIn struct {
	Fingerprint    string `json:fingerPrint`
	curTime        string `json:curTime`
	CustomerName   string `json:customerName`
	CustomerSecret string `json:customerSecret`
}

type licenseStatus struct {
	FeatureName string `json:featureName`
	MaxUsage    string `json:maxUsage`
	CurUsage    string `json:curUsage`
	LicenseHash string `json:licenseHash`
}

type deviceHBPunchOut struct {
	Status       int             `json:statusCode`
	StatusString string          `json:statusString`
	curTime      string          `json:curTime`
	CurLicStatus []licenseStatus `json:currentLicenseStatus`
}

type billingAddPurchaseReq struct {
	CustomerName string `json:customerName`
	FeatureName  string `json:featureName`
	LicenseCount int    `json:licenseCount`
	UsagePeriod  int    `json:usagePeriod` //in seconds
}

type billingAddPurchaseResp struct {
	Status       int    `json:statusCode`
	StatusString string `json:statusString`
}

type uiMarkDeviceRMAReq struct {
	CustomerName   string `json:customerName`
	CustomerSecret string `json:customerSecret`
}

type uiMarkDeviceRMAResp struct {
	Status       int    `json:statusCode`
	StatusString string `json:statusString`
}

type uiCustomerReportReq struct {
	CustomerName   string `json:customerName`
	CustomerSecret string `json:customerSecret`
}

type uiCustomerReportResp struct {
	Status       int             `json:statusCode`
	StatusString string          `json:statusString`
	CurLicStatus []licenseStatus `json:currentLicenseStatus`
}
