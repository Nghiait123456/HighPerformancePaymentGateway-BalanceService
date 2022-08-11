

request check balance ==>
                        mu_lock for one partner ==>
                                                    balance check in local ram ==>
                        free_lock
                       ==>  insert to db
                                        => if success : response  success

                                        ==> if fail:  ==> get Mu_lock
                                                                    ==> roolback balance
                                                      ==> free MU Lock
                                                      ==> response fail
