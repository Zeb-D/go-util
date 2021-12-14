package dubbo

import (
	"encoding/base64"
	"fmt"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/stretchr/testify/assert"
	"testing"
)

const DATA = "lE0XamF2YS51dGlsLkxpbmtlZEhhc2hNYXAHc3VjY2Vzc1QGc3RhdHVzAm9rCWVycm9yQ29kZU4IZXJyb3JNc2cAAXRMAAABdkV5bFwGcmVzdWx0SARjb2RlTgRkYXRhSA5zaG9wQ2FydE1vZGVsc3ETamF2YS51dGlsLkFycmF5TGlzdEgKY2FydE1vZGVsc3ORSA5pbnZhbGlkTWVzc2FnZU4McG9zdGFnZVByaWNlTgNudW2RC25lZWRDaGFuZ2VkRgttYWluUGljdHVyZTAjZWNvbW1lcmNlLzE1ODMzMjMzNjI5Y2E1YjgzYjI5Ni5qcGcKZ2lmdEFjdGl2ZU4FdGl0bGUQQXV0b1Rlc3Q4MDA0NDg2Mg5hY3Rpdml0eU1vZGVsc04FcHJpY2X4ZA1jb21tb2RpdHlDb2RlDkNNOXF4cnU4czRrZm5oD2F0dHJWYWx1ZU1vZGVsc3GRSBBjb21tb2RpdHlTa3VDb2RlDlNLOXF4cnU4c2gxeGVsC2dtdE1vZGlmaWVkTghhdHRyQ29kZQx4anRoMzViR1M5SU0LYXR0clZhbHVlSWT77Qtjb21tb2RpdHlJZFkADdPoCWdtdENyZWF0ZU4EdW5pdE4GYXR0cklk4w5jb21tb2RpdHlTa3VJZFkADmR1DWNvbW1vZGl0eUNvZGUOQ005cXhydThzNGtmbmgJYXR0clZhbHVlA+ebtOa1geeUtQJpZFkADbStBWNsYXNzMDljb20udHV5YS5saW9uLmNsaWVudC5za3UubW9kZWwuQ29tbW9kaXR5U2t1QXR0clZhbHVlTW9kZWwIYXR0ck5hbWVOBnN0YXR1c05aBmRldGFpbEgPcG9zdGFnZUluZm9MaXN0TgtnbXRNb2RpZmllZE4Hc2t1TGlzdHGRSAtnbXRNb2RpZmllZE4FcHJpY2XIZA9za3VBdHRyVmFsdWVNYXBIDHhqdGgzNWJHUzlJTU0fY29tLmFsaWJhYmEuZmFzdGpzb24uSlNPTk9iamVjdAhhdHRyQ29kZQx4anRoMzViR1M5SU0LYXR0clZhbHVlSWT77QlhdHRyVmFsdWUD55u05rWB55S1WloNY29tbW9kaXR5Q29kZQ5DTTlxeHJ1OHM0a2ZuaAJpZE4Fc3RvY2uaCWdtdENyZWF0ZU4FY2xhc3MwPmNvbS50dXlhLnNoZW5zaG91LmNsaWVudC5tb2RlbC5hcHAuY29tbW9kaXR5Lm1vZGVsLkFwcFNrdU1vZGVsC3NvdXJjZVByaWNlyGQHc2t1Q29kZQ5TSzlxeHJ1OHNoMXhlbFoWcmVjb21tZW5kQ29tbW9kaXR5TGlzdE4MdmlkZW9QaWN0dXJlTg5tYXhTb3VyY2VQcmljZchkBXZpZGVvTgV0aXRsZRBBdXRvVGVzdDgwMDQ0ODYyCnNlbGxlckNvZGULdGVzdF9zZWxsZXIQc2t1QXR0ck1vZGVsTGlzdHGRSAtnbXRNb2RpZmllZE4IYXR0ckNvZGUMeGp0aDM1YkdTOUlNDWF0dHJHcm91cENvZGUDU0tVC2F0dHJHcm91cElk4g1yZWxldmFuY2VMaXN0TghsYW5nTGlzdE4JZ210Q3JlYXRlTghyZXF1aXJlZFQIYXR0clR5cGWQBHVuaXQBaAdkZWxldGVkRgJpZOMSYXR0clZhbHVlTW9kZWxMaXN0cZFIDWF0dHJWYWx1ZURlc2NOC2dtdE1vZGlmaWVkTghhdHRyQ29kZQx4anRoMzViR1M5SU0GYXR0cklk4wlhdHRyVmFsdWUD55u05rWB55S1Amlk++0IbGFuZ0xpc3ROCWdtdENyZWF0ZU4FY2xhc3MwLmNvbS50dXlhLmxpb24uY2xpZW50LmF0dHIubW9kZWwuQXR0clZhbHVlTW9kZWwFb3JkZXJOWgVjbGFzczApY29tLnR1eWEubGlvbi5jbGllbnQuYXR0ci5tb2RlbC5BdHRyTW9kZWwIYXR0ck5hbWUE5L6b55S15pa55byPWgxpbnN0YWxsUHJpY2WQCHN1YlRpdGxlABBjb21tb2RpdHlEZXRhaWxzcZEwI2Vjb21tZXJjZS8xNTgzMzIzMzkwZjIxNzE3YzI0Y2QuanBnCG1hbGxDb2RlCXRlc3RfY29kZQ1jb21tb2RpdHlDb2RlDkNNOXF4cnU4czRrZm5oC3Bvc3RhZ2VUeXBlkQJpZE4FY2xhc3MwRGNvbS50dXlhLnNoZW5zaG91LmNsaWVudC5tb2RlbC5hcHAuY29tbW9kaXR5Lm1vZGVsLkFwcENvbW1vZGl0eU1vZGVsC3BpY3R1cmVMaXN0cZEwI2Vjb21tZXJjZS8xNTgzMzIzMzYyOWNhNWI4M2IyOTYuanBnCnN0YXJ0UHJpY2VOCXBvc3RhZ2VJZE4McG9zdGFnZVByaWNlyGQIbWFsbFBhdGgACmdpZnRBY3RpdmVOEGZyZWVMb2NhdGlvbkxpc3RODnNlcnZpY2VEZXNjVXJpTglnbXRDcmVhdGVOC2luc3RhbGxUeXBlkQhtaW5QcmljZchkB3B1Ymxpc2iRDm1pblNvdXJjZVByaWNlyGQKdG90YWxTdG9ja5oIbWF4UHJpY2XIZAdjb2xsZWN0RgZzdGF0dXNOWgVzdG9ja5oFY2xhc3MwRGNvbS50dXlhLnNoZW5zaG91LmNsaWVudC5tb2RlbC5yZXRhaWwuY2FydC5tb2RlbC5SZXRhaWxDYXJ0SXRlbU1vZGVsB3NrdUNvZGUOU0s5cXhydThzaDF4ZWwGc3RhdHVzkVpIDmludmFsaWRNZXNzYWdlTgxwb3N0YWdlUHJpY2VOA251bZILbmVlZENoYW5nZWRGC21haW5QaWN0dXJlMCRlY29tbWVyY2UvMTU5OTQ3NjUxNzMxZDJhOWUzZmQzLmpwZWcKZ2lmdEFjdGl2ZU4FdGl0bGUwIOa1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlea1i+ivlQ5hY3Rpdml0eU1vZGVsc04FcHJpY2XhDWNvbW1vZGl0eUNvZGUOQ005d2Q0aDB0MW44bHIPYXR0clZhbHVlTW9kZWxzcZFIEGNvbW1vZGl0eVNrdUNvZGUOU0s5d2Q0aDExMHRrY3kLZ210TW9kaWZpZWROCGF0dHJDb2RlDDRuT2VBV1lhVk45YwthdHRyVmFsdWVJZPvyC2NvbW1vZGl0eUlkWQAQAasJZ210Q3JlYXRlTgR1bml0TgZhdHRySWT77A5jb21tb2RpdHlTa3VJZFkAEKXADWNvbW1vZGl0eUNvZGUOQ005d2Q0aDB0MW44bHIJYXR0clZhbHVlClNLVV9WQUxVRTICaWRZAA/O6AVjbGFzczA5Y29tLnR1eWEubGlvbi5jbGllbnQuc2t1Lm1vZGVsLkNvbW1vZGl0eVNrdUF0dHJWYWx1ZU1vZGVsCGF0dHJOYW1lTgZzdGF0dXNOWgZkZXRhaWxID3Bvc3RhZ2VJbmZvTGlzdE4LZ210TW9kaWZpZWROB3NrdUxpc3RykUgLZ210TW9kaWZpZWROBXByaWNlkQ9za3VBdHRyVmFsdWVNYXBIDDRuT2VBV1lhVk45Y02SCGF0dHJDb2RlDDRuT2VBV1lhVk45YwthdHRyVmFsdWVJZPvyCWF0dHJWYWx1ZQpTS1VfVkFMVUUyWloNY29tbW9kaXR5Q29kZQ5DTTl3ZDRoMHQxbjhscgJpZE4Fc3RvY2vISglnbXRDcmVhdGVOBWNsYXNzMD5jb20udHV5YS5zaGVuc2hvdS5jbGllbnQubW9kZWwuYXBwLmNvbW1vZGl0eS5tb2RlbC5BcHBTa3VNb2RlbAtzb3VyY2VQcmljZchkB3NrdUNvZGUOU0s5d2Q0aDExMHRrY3laSAtnbXRNb2RpZmllZE4FcHJpY2WRD3NrdUF0dHJWYWx1ZU1hcEgMNG5PZUFXWWFWTjljTZIIYXR0ckNvZGUMNG5PZUFXWWFWTjljC2F0dHJWYWx1ZUlk+/EJYXR0clZhbHVlClNLVV9WQUxVRTFaWg1jb21tb2RpdHlDb2RlDkNNOXdkNGgwdDFuOGxyAmlkTgVzdG9ja8hjCWdtdENyZWF0ZU4FY2xhc3MwPmNvbS50dXlhLnNoZW5zaG91LmNsaWVudC5tb2RlbC5hcHAuY29tbW9kaXR5Lm1vZGVsLkFwcFNrdU1vZGVsC3NvdXJjZVByaWNlyGQHc2t1Q29kZQ5TS2EwbXMxNTNkeTZ2cloWcmVjb21tZW5kQ29tbW9kaXR5TGlzdE4MdmlkZW9QaWN0dXJlTg5tYXhTb3VyY2VQcmljZchkBXZpZGVvTgV0aXRsZTAg5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+V5rWL6K+VCnNlbGxlckNvZGULdGVzdF9zZWxsZXIQc2t1QXR0ck1vZGVsTGlzdHGRSAtnbXRNb2RpZmllZE4IYXR0ckNvZGUMNG5PZUFXWWFWTjljDWF0dHJHcm91cENvZGUDU0tVC2F0dHJHcm91cElk4g1yZWxldmFuY2VMaXN0TghsYW5nTGlzdE4JZ210Q3JlYXRlTghyZXF1aXJlZFQIYXR0clR5cGWRBHVuaXQCY20HZGVsZXRlZEYCaWT77BJhdHRyVmFsdWVNb2RlbExpc3RykUgNYXR0clZhbHVlRGVzY04LZ210TW9kaWZpZWROCGF0dHJDb2RlDDRuT2VBV1lhVk45YwZhdHRySWT77AlhdHRyVmFsdWUKU0tVX1ZBTFVFMgJpZPvyCGxhbmdMaXN0TglnbXRDcmVhdGVOBWNsYXNzMC5jb20udHV5YS5saW9uLmNsaWVudC5hdHRyLm1vZGVsLkF0dHJWYWx1ZU1vZGVsBW9yZGVyTlpIDWF0dHJWYWx1ZURlc2NOC2dtdE1vZGlmaWVkTghhdHRyQ29kZQw0bk9lQVdZYVZOOWMGYXR0cklk++wJYXR0clZhbHVlClNLVV9WQUxVRTECaWT78QhsYW5nTGlzdE4JZ210Q3JlYXRlTgVjbGFzczAuY29tLnR1eWEubGlvbi5jbGllbnQuYXR0ci5tb2RlbC5BdHRyVmFsdWVNb2RlbAVvcmRlck5aBWNsYXNzMCljb20udHV5YS5saW9uLmNsaWVudC5hdHRyLm1vZGVsLkF0dHJNb2RlbAhhdHRyTmFtZQfku4HotLXmtYvor5VTS1VaDGluc3RhbGxQcmljZZAIc3ViVGl0bGUAEGNvbW1vZGl0eURldGFpbHNxkTAkZWNvbW1lcmNlLzE1OTk0NzY1Njc1ZmFmZDQwNWFjZS5qcGVnCG1hbGxDb2RlCXRlc3RfY29kZQ1jb21tb2RpdHlDb2RlDkNNOXdkNGgwdDFuOGxyC3Bvc3RhZ2VUeXBlkQJpZE4FY2xhc3MwRGNvbS50dXlhLnNoZW5zaG91LmNsaWVudC5tb2RlbC5hcHAuY29tbW9kaXR5Lm1vZGVsLkFwcENvbW1vZGl0eU1vZGVsC3BpY3R1cmVMaXN0cZEwJGVjb21tZXJjZS8xNTk5NDc2NTE3MzFkMmE5ZTNmZDMuanBlZwpzdGFydFByaWNlTglwb3N0YWdlSWRODHBvc3RhZ2VQcmljZchkCG1hbGxQYXRoAApnaWZ0QWN0aXZlThBmcmVlTG9jYXRpb25MaXN0Tg5zZXJ2aWNlRGVzY1VyaU4JZ210Q3JlYXRlTgtpbnN0YWxsVHlwZZEIbWluUHJpY2WRB3B1Ymxpc2iRDm1pblNvdXJjZVByaWNlyGQKdG90YWxTdG9ja8itCG1heFByaWNlkQdjb2xsZWN0RgZzdGF0dXNOWgVzdG9ja8hKBWNsYXNzMERjb20udHV5YS5zaGVuc2hvdS5jbGllbnQubW9kZWwucmV0YWlsLmNhcnQubW9kZWwuUmV0YWlsQ2FydEl0ZW1Nb2RlbAdza3VDb2RlDlNLOXdkNGgxMTB0a2N5BnN0YXR1c5FaSA5pbnZhbGlkTWVzc2FnZU4McG9zdGFnZVByaWNlTgNudW2SC25lZWRDaGFuZ2VkVAttYWluUGljdHVyZTAjZWNvbW1lcmNlLzE1ODMzMjMzNjI5Y2E1YjgzYjI5Ni5qcGcKZ2lmdEFjdGl2ZU4FdGl0bGUW5o6l5Y+j5rWL6K+VMTYwNjE5NzY2MjgwNS45MzUzDmFjdGl2aXR5TW9kZWxzTgVwcmljZTwnEA1jb21tb2RpdHlDb2RlDkNNOXF2cThob3gxNHkzD2F0dHJWYWx1ZU1vZGVsc3GRSBBjb21tb2RpdHlTa3VDb2RlDlNLOXF2cThocDlpbW1oC2dtdE1vZGlmaWVkTghhdHRyQ29kZQx4anRoMzViR1M5SU0LYXR0clZhbHVlSWT77Qtjb21tb2RpdHlJZFkADZ2xCWdtdENyZWF0ZU4EdW5pdE4GYXR0cklk4w5jb21tb2RpdHlTa3VJZFkADi5YDWNvbW1vZGl0eUNvZGUOQ005cXZxOGhveDE0eTMJYXR0clZhbHVlA+ebtOa1geeUtQJpZFkADYJ0BWNsYXNzMDljb20udHV5YS5saW9uLmNsaWVudC5za3UubW9kZWwuQ29tbW9kaXR5U2t1QXR0clZhbHVlTW9kZWwIYXR0ck5hbWVOBnN0YXR1c05aBmRldGFpbEgPcG9zdGFnZUluZm9MaXN0TgtnbXRNb2RpZmllZE4Hc2t1TGlzdHGRSAtnbXRNb2RpZmllZE4FcHJpY2XUJxAPc2t1QXR0clZhbHVlTWFwSAx4anRoMzViR1M5SU1NkghhdHRyQ29kZQx4anRoMzViR1M5SU0LYXR0clZhbHVlSWT77QlhdHRyVmFsdWUD55u05rWB55S1WloNY29tbW9kaXR5Q29kZQ5DTTlxdnE4aG94MTR5MwJpZE4Fc3RvY2uRCWdtdENyZWF0ZU4FY2xhc3MwPmNvbS50dXlhLnNoZW5zaG91LmNsaWVudC5tb2RlbC5hcHAuY29tbW9kaXR5Lm1vZGVsLkFwcFNrdU1vZGVsC3NvdXJjZVByaWNl1CcQB3NrdUNvZGUOU0s5cXZxOGhwOWltbWhaFnJlY29tbWVuZENvbW1vZGl0eUxpc3RODHZpZGVvUGljdHVyZU4ObWF4U291cmNlUHJpY2XUJxAFdmlkZW9OBXRpdGxlFuaOpeWPo+a1i+ivlTE2MDYxOTc2NjI4MDUuOTM1MwpzZWxsZXJDb2RlC3Rlc3Rfc2VsbGVyEHNrdUF0dHJNb2RlbExpc3RxkUgLZ210TW9kaWZpZWROCGF0dHJDb2RlDHhqdGgzNWJHUzlJTQ1hdHRyR3JvdXBDb2RlA1NLVQthdHRyR3JvdXBJZOINcmVsZXZhbmNlTGlzdE4IbGFuZ0xpc3ROCWdtdENyZWF0ZU4IcmVxdWlyZWRUCGF0dHJUeXBlkAR1bml0AWgHZGVsZXRlZEYCaWTjEmF0dHJWYWx1ZU1vZGVsTGlzdHGRSA1hdHRyVmFsdWVEZXNjTgtnbXRNb2RpZmllZE4IYXR0ckNvZGUMeGp0aDM1YkdTOUlNBmF0dHJJZOMJYXR0clZhbHVlA+ebtOa1geeUtQJpZPvtCGxhbmdMaXN0TglnbXRDcmVhdGVOBWNsYXNzMC5jb20udHV5YS5saW9uLmNsaWVudC5hdHRyLm1vZGVsLkF0dHJWYWx1ZU1vZGVsBW9yZGVyTloFY2xhc3MwKWNvbS50dXlhLmxpb24uY2xpZW50LmF0dHIubW9kZWwuQXR0ck1vZGVsCGF0dHJOYW1lBOS+m+eUteaWueW8j1oMaW5zdGFsbFByaWNlkAhzdWJUaXRsZQAQY29tbW9kaXR5RGV0YWlsc3GRMCNlY29tbWVyY2UvMTU4MzMyMzM5MGYyMTcxN2MyNGNkLmpwZwhtYWxsQ29kZQl0ZXN0X2NvZGUNY29tbW9kaXR5Q29kZQ5DTTlxdnE4aG94MTR5Mwtwb3N0YWdlVHlwZZECaWROBWNsYXNzMERjb20udHV5YS5zaGVuc2hvdS5jbGllbnQubW9kZWwuYXBwLmNvbW1vZGl0eS5tb2RlbC5BcHBDb21tb2RpdHlNb2RlbAtwaWN0dXJlTGlzdHGRMCNlY29tbWVyY2UvMTU4MzMyMzM2MjljYTViODNiMjk2LmpwZwpzdGFydFByaWNlTglwb3N0YWdlSWRODHBvc3RhZ2VQcmljZchkCG1hbGxQYXRoAApnaWZ0QWN0aXZlThBmcmVlTG9jYXRpb25MaXN0Tg5zZXJ2aWNlRGVzY1VyaU4JZ210Q3JlYXRlTgtpbnN0YWxsVHlwZZEIbWluUHJpY2XUJxAHcHVibGlzaJEObWluU291cmNlUHJpY2XUJxAKdG90YWxTdG9ja5EIbWF4UHJpY2XUJxAHY29sbGVjdEYGc3RhdHVzTloFc3RvY2uRBWNsYXNzMERjb20udHV5YS5zaGVuc2hvdS5jbGllbnQubW9kZWwucmV0YWlsLmNhcnQubW9kZWwuUmV0YWlsQ2FydEl0ZW1Nb2RlbAdza3VDb2RlDlNLOXF2cThocDlpbW1oBnN0YXR1c5FaCHNob3BJY29uMDJlY29tbWVyY2UvMTU1ODY2NzA3NTE3NF8xNTU4NjY3MDc1X3E1YW1ibWEyODdqLnBuZwhzaG9wTmFtZQTlrZDlurjllYbln44FY2xhc3MwRGNvbS50dXlhLnNoZW5zaG91LmNsaWVudC5tb2RlbC5yZXRhaWwuY2FydC5tb2RlbC5SZXRhaWxTaG9wQ2FydE1vZGVsWghtYXhMaW1pdE4FY291bnSUC2ludmFsaWRMaXN0cZFIDmludmFsaWRNZXNzYWdlTgxwb3N0YWdlUHJpY2VOA251bZELbmVlZENoYW5nZWRGC21haW5QaWN0dXJlMCNlY29tbWVyY2UvMTU4MzMyMzM2MjljYTViODNiMjk2LmpwZwpnaWZ0QWN0aXZlTgV0aXRsZRjmjqXlj6PmtYvor5XnvJbovpExNjA1MTgxNjY4NzU0LjUyNTQOYWN0aXZpdHlNb2RlbHNOBXByaWNlPBPsDWNvbW1vZGl0eUNvZGUOQ01hMm96Z2N6Z2pjNWUPYXR0clZhbHVlTW9kZWxzTgZkZXRhaWxOBXN0b2NrmQVjbGFzczBEY29tLnR1eWEuc2hlbnNob3UuY2xpZW50Lm1vZGVsLnJldGFpbC5jYXJ0Lm1vZGVsLlJldGFpbENhcnRJdGVtTW9kZWwHc2t1Q29kZQ5TS2Eyb3pnZDAzMGY5dwZzdGF0dXOQWgVjbGFzczBAY29tLnR1eWEuc2hlbnNob3UuY2xpZW50Lm1vZGVsLnJldGFpbC5jYXJ0Lm1vZGVsLlJldGFpbENhcnRNb2RlbFoBdEwAAAF2RXlsXAdzdWNjZXNzVAdtZXNzYWdlTgVjbGFzczAzY29tLnR1eWEuY29tbWVyY2lhbC5hcGkubW9kZWxzLmNvbW1vbi5TZXJ2aWNlUmVzdWx0WgN0YWcMU0hFTkNMSTpSRVNQWkgGb3V0cHV0AjUyBWR1YmJvBTIuMC4yWg=="

func TestDubboHessian(t *testing.T) {
	bytes, err := base64.StdEncoding.DecodeString(DATA)
	assert.Nil(t, err)
	fmt.Println(bytes)

	// demo a decoder to decode buffer from client
	d := hessian.NewDecoder(bytes)
	lengthObj, _ := d.Decode()
	assert.True(t, lengthObj == int32(4))
	length := lengthObj.(int32)
	fmt.Println(length)
	fmt.Println("resp")
	o, err := d.DecodeValue()
	assert.Nil(t, err)
	fmt.Println(o)

	o, err = d.DecodeValue()
	assert.Nil(t, err)
	fmt.Println("att")
	fmt.Println(o)
}