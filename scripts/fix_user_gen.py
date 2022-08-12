import string
import pandas as pd

a = range(2000)
b = []
name_prefix = "stress_test_"
start_num = 10000000
for i in a :
    b.append(name_prefix + str(start_num + i))

data_frame = pd.DataFrame({"username": b})

data_frame.to_csv("fix_2000.csv", index=False, sep=",")