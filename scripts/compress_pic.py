import os
from PIL import Image

DIR = "/Users/yuan.ding/Desktop/code/entry_task/images/test_small_resized.jpg"

img = Image.open(DIR)

width, height = img.size[0], img.size[1]

newPath = "/Users/yuan.ding/Desktop/code/entry_task/images/test_small_resized.jpg"
new_width = int(width / 2)
new_height = int(height / 2)

new_img = img.resize((new_width, new_height))
new_img.save(newPath)