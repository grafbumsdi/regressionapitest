# regressionapitest
API calls done via google go application

Sample calls:
```bash
regressionapitest.exe -serveraddress 192.168.11.25 -apicalls api/v1/wikifolios/AFSDFWEFWEF,aaaaaaaaaaapi/v1/wikifolios/AFSDFWEFWEF
```
```bash
regressionapitest.exe -serveraddress 192.168.11.43 -loglevel Trace -logfile apitests_43.log
```
```bash
regressionapitest.exe -serveraddress 192.168.11.23/api/v1 -apicalls wikifolios,trades,import/wikifolios
```