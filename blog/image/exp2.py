from matplotlib import image #for working with images
from matplotlib import pyplot #for plotting stuff on graphical gui

matrix = image.imread('p.png') #read file
#indexes for rgb values in a array of matrix values 
red = 0
green = 1
blue = 2
opacity = 3 #PNG format allows controlling transparency of image

for row in range(len(matrix)): #cycle through each row in the image
    for column in range(len(matrix[row])): #cycle through column of pixel in the image
        for value in range(len(matrix[row][column])): #cycle through rgb values within a pixel in the image
            # NEGATIVE
            if(value!=opacity): #need not change opacity
                matrix[row][column][value]=1-matrix[row][column][value] #negate rgb values (max is 1) so that dark becomes light and vv
            
            # ADD A VALUE
            # if(value!=opacity): #need not change opacity
                # matrix[row][column][value]=matrix[row][column][value]+0.5 # 0 is least 1 is max

            # SUBTRACT A VALUE
            # if(value!=opacity): #need not change opacity
                # matrix[row][column][value]=matrix[row][column][value]-0.5 # 0 is least 1 is max

            # MULTIPLY A VALUE
            # if(value!=opacity): #need not change opacity
                # matrix[row][column][value]=matrix[row][column][value]*2 # 0 is least 1 is max

            # DIVIDE A VALUE
            # if(value!=opacity): #need not change opacity
                # matrix[row][column][value]=matrix[row][column][value]/2 # 0 is least 1 is max
            
            # Darker blacks, brighter whites
            # if(value!=opacity):
            #     if(matrix[row][column][value]<0.5):
            #         matrix[row][column][value]=matrix[row][column][value]*0.5
            #     else:
            #         matrix[row][column][value]=matrix[row][column][value]*1.5

pyplot.axis('off')
pyplot.imshow(matrix)
pyplot.show()