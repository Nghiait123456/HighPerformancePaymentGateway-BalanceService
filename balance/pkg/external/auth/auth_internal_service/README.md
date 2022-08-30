/**
auth internal in microservice: high performance + security + handle all race conditions + zero down time when update new
secret
process:
you have:    secret manager remote(encry AllData) save all secret key
all service ==> cal to secret manager remote ==> get secret ==> auth ||||||||||||||||||||||||||| secret manager remote: ==> auto change secret for time

*/

use please view example