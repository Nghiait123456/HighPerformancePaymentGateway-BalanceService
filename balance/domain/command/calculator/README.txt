

request check balance ==>
                        mu_lock for one partner
                                                    ==>  balance check in local ram
                                                    ==>  insert to db
                        free_lock

                        => if success : response  success
                        ==> if fail:  response  fail
