import time as t #using it for sleep

#speed of track
def setTempo(time):
    Clock.bpm=time #set tempo(its like speed of beats but not exactly speed)
    Clock.mod(16)

# guitar sounds and tracks
def guitar(fill1,fill2):
    if(fill1):
        #for slow guitar play in the music, with some echo
        p1>>pluck([0,0,0, #these nos are for keys, next number for next key viz. added on new lines for better understanding of keys
            2,
            5,5,5,
            4,
            2,2,2,
            2,
            0,
            -1,-1,
            -1,0],dur=[1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1/2,1/2],sus=1) 
    if(fill2): 
        #for the second overlapping guitar with double notes
        p2>>pluck([0,0,0,0,0,0,
          2,2,2,2,
          5,5,5,5,5,5,
          4,4,
          2,2,2,2,2,2,
          2,2,2,2,
          0,
          -1,-1,-1,-1,-1,-1,
          -1,0
          ],
          dur=[1/2,1/2,1/2,1/2,1/2,1/2,
              1/4,1/4,1/4,1/4,
              1/2,1/2,1/2,1/2,1/2,1/2,
              1/2,1/2,
              1/2,1/2,1/2,1/2,1/2,1/2,
              1/4,1/4,1/4,1/4,
              1/2,
              1/2,1/2,1/2,1/2,1/2,1/2,
              1/4,1/4
              ],amp=0.5,sus=3)

#beats for track
def regularBeat(preBeatDrop=2,trigger=True,ampp=4):
    if(trigger):
        p3>>play("x",dur=[preBeatDrop],amp=ampp)
    else:
        p3.stop()
    
#beats change before DROP    
def beforeDropBeat():
    regularBeat(1)
    t.sleep(16)
    regularBeat(1/2)
    t.sleep(8)
    regularBeat(1/4)
    t.sleep(6)
    regularBeat(1/8,True,4.5)
    t.sleep(4)
    regularBeat(0,True)

#track after the drop
def drop():
    guitar(1,0)
    p4>>play("x-",dur=[1/2,1/2],sample=[17,17],sus=5,amp=[4,1])

#stop music
def stop():
    Clock.clear()
    
#all components of song
def alanWalkerFaded():
    setTempo(80)
    guitar(1,0)
    t.sleep(16)
    guitar(1,1)
    t.sleep(8)
    regularBeat()
    t.sleep(2)
    beforeDropBeat()
    drop()
    t.sleep(40)
    stop()

#play the song !
alanWalkerFaded()

