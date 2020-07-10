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
        # RED ONLY
        matrix[row][column][red]=matrix[row][column][red]
        matrix[row][column][green]=0 
        matrix[row][column][blue]=0

        # GREEN ONLY
        # matrix[row][column][red]=0
        # matrix[row][column][green]=matrix[row][column][green] 
        # matrix[row][column][blue]=0

        # BLUE ONLY
        # matrix[row][column][red]=0
        # matrix[row][column][green]=0
        # matrix[row][column][blue]=matrix[row][column][blue]

        # RED GREEN ONLY
        # matrix[row][column][red]=matrix[row][column][red] 
        # matrix[row][column][green]=matrix[row][column][green] 
        # matrix[row][column][blue]=0

        # RED BLUE ONLY
        # matrix[row][column][red]=matrix[row][column][red] 
        # matrix[row][column][green]=0
        # matrix[row][column][blue]=matrix[row][column][blue] 

        # GREEN BLUE ONLY
        # matrix[row][column][red]=matrix[row][column][red] 
        # matrix[row][column][green]=matrix[row][column][green] 
        # matrix[row][column][blue]=0

        # RED GREEN BLUE 50% TRANSPARENT
        # matrix[row][column][opacity]=0.5
        # matrix[row][column][opacity]=0.5 
        # matrix[row][column][opacity]=0.5

pyplot.axis('off')
pyplot.imshow(matrix)
pyplot.show()