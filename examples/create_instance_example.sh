# first we create instance
../tcactl create cnf mwcapp02 mwcapp02-instance01
watch -c ../tcactl get cnfi mwcapp02-instance

watch -c ../tcactl update cnf reconfigure mwcapp02-instance sample_replica.yaml

##
../tcactl update cnf down mwcapp02-instance01 --pool dallas
##
../tcactl update cnf up mwcapp02-instance01 --pool atlanta

