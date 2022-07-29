I was use : https://github.com/Nghiait123456/backendForLoadtest . I was use c5.24x.larger for backend, c5.24x.larger for
loadtest. </br>
Result from result.txt: </br>
+) with one instance loadtest, i have 85000 rps </br>
+) with four instance loadtest, i have 210 000 to 230 000 rqs </br>, cpu ~ 30 % </br>
+) i don't have enough Ec2 resources, i consult some benmark sources, with my close spec, 8 instnace loadtest, can
achieve ~500K rqs, cpu ~ 65% </br>

===>     Fiber is frame has an impressive http load capacity. I'll dig into this later, but I'll choose it for my
service balance. </br>

