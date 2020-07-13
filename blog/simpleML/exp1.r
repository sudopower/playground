import pandas as pd #to work with CSVs
mydata_path = '../input/simple.csv' #import file
mydata = pd.read_csv(mydata_path) #read the file

mydata.head() #first few lines

#The model will look and a, b & c
features=['a','b','c']
X = mydata[features]

#The model will then look at the output y
y = mydata.y

from sklearn.tree import DecisionTreeRegressor #Decision tree model ready code
my_model = DecisionTreeRegressor() # initialize an empty model
my_model.fit(X, y) #make a decision tree that fits our data

#We create 5 test conditions
case_one =[1,1,1] 
case_two =[1,0,1] 
case_three =[0,0,0] 
case_four =[1,0,0] 
case_five =[1,0,1] 

y_test=[case_one,case_two,case_three,case_four,case_five]
print(my_model.predict(y_test))

#OUTPUT for 5 predictions
#[1. 1. 0. 1. 1.]


new_case =[0,0,1]
y_test=[new_case]
print(my_model.predict(y_test))

#OUTPUT
#[0.]
