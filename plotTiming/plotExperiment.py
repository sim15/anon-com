"""
=============================================
Generate polygons to fill under 3D line graph
=============================================

Demonstrate how to create polygons which fill the space under a line
graph. In this example polygons are semi-transparent, creating a sort
of 'jagged stained glass' effect.
"""

# from ensurepip import version
# from cProfile import label
# from mpl_toolkits.mplot3d import Axes3D
# from matplotlib.collections import PolyCollection
# import matplotlib.pyplot as plt
# from matplotlib import colors as mcolors
# import numpy as np
import json
import math
import statistics
# import matplotlib.ticker as ticker


# fig = plt.figure()
# ax = fig.gca(projection='3d')


# def cc(arg):
#     return mcolors.to_rgba(arg, alpha=0.6)

# xs = np.arange(0, 10, 0.4)
# verts = []
# zs = [0.0, 1.0, 2.0, 3.0]
# for z in zs:
#     ys = np.random.rand(len(xs))
#     ys[0], ys[-1] = 0, 0
#     verts.append(list(zip(xs, ys)))

# print(xs)
# print(len(verts), len(verts[0]))
# print(zs)

valuesTime = {}


with open('experiment1.json') as f:
    data = json.load(f)

    for experiment in data:
        n = int(math.log2(experiment["NumBoxes"]))
        if n > 0:
            l = experiment["message_length"]
            t = sum(statistics.mean(i) for i in [experiment["construction1_ms"][0]] ) / 1000
            if not (l in valuesTime):
                valuesTime[l] = []  
        
            valuesTime[l].append(t)
    
    # for l in valuesTime.keys():
    #     valuesTime[l].append((valuesTime[l][-1][0]+1, 0))
    
        
verts = list(valuesTime.values())
zs = [key for key in valuesTime]
# print(len(verts), len(verts[0]))
# print(zs)

print(verts)
print(zs)

# poly = PolyCollection(verts,
#  linewidths=[2 for i in range(len(verts))],
#  closed=False,
#   edgecolors='blue',
#    facecolor=["none" for i in range(len(verts))])
# # poly.set_alpha(0.2)
# ax.add_collection3d(poly, zs=zs, zdir='y')

# ax.set_xlabel('# of mailboxes')
# ax.set_xlim3d(8, 20)
# # ax.set_xscale("log")
# ax.set_xticks(np.arange(8, 21, 4))
# ax.xaxis.set_major_formatter(ticker.FuncFormatter(
#         lambda v, x: "$2^{%i}$" % v))


# ax.set_title('Total Server-side Processing (writing + authentication)')
# ax.set_ylabel('message size (B)')
# ax.set_ylim3d(0, 1050)

# ax.set_yticks([100, 300, 500, 750, 1000])

# ax.set_zlabel('time (s)')
# ax.set_zlim3d(0, 6)

# # ax.scatter([0, 20],[100, 100],[0,0], label="endpoints")

# ax.view_init(elev=25, azim=-130)

# plt.show()