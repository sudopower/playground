from PIL import Image, ImageDraw, ImageFont #for drawing image
from matplotlib import image #can be done with PIL probably, just for code reuse

matrix = image.imread('apple.png') #read file

#indexes for rgb values in a array of matrix values 
red = 0
green = 1
blue = 2
opacity = 3 #PNG format allows controlling transparency of image

scaling = 20 #scale factor for text size, spacing, line heighting, image dimensions
character = "\U0001F34E" #red apple emoticon unicode
canvasColor = (255, 255, 255)
dimension = len(matrix) # height = width, sq image

font_path = "seguiemj.ttf" #font to work with emoticons
font_size = scaling #set font size

img = Image.new('RGB', (dimension*scaling, dimension*scaling), color = canvasColor) #create canvas for image
d = ImageDraw.Draw(img) # draw into this canvas
font = ImageFont.truetype(font_path, font_size) #set font properties

for row in range(dimension): #cycle through each row in the image
    spaceCounter=0 # to space out the pixel apples
    for column in range(len(matrix[row])): #cycle through column of pixel in the image
        # color as per value
        redVal=matrix[row][column][red] #redness
        greenVal=matrix[row][column][green] #greeness
        blueVal=matrix[row][column][blue] #blueness
        hText = column+spaceCounter*scaling #scale and space out every pixel apple
        vText = row*scaling #space out every row
        d.text((hText,vText), character, fill=(int(redVal*255), int(greenVal*255), int(blueVal*255)),font=font) #place the text with color as per original image in corresponding coordinates
        spaceCounter+=1 #push every pixel apple in row further away
 
img.save('apples.png')