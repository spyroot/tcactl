# first we create instance
../tcactl create cnf mwcapp02 mwcapp02-test01 --block
#watch -c ../tcactl get cnfi mwcapp02-instance

#../tcactl update cnf reconfigure mwcapp02-instance sample_replica.yaml
../tcactl delete cnf mwcapp02-test01 --force

