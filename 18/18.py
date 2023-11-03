import json

import matplotlib.pyplot as plt
import numpy as np


def explode(data):
    size = np.array(data.shape) * 2
    data_e = np.zeros(size - 1, dtype=data.dtype)
    data_e[::2, ::2, ::2] = data
    return data_e


for f in ["ex", "in"]:
    with open(f'{f}.txt') as fh:
        cubes = []
        xs, ys, zs = [], [], []
        for line in fh.readlines():
            if len(line) == 0:
                continue
            x, y, z = map(int, line.split(","))
            cubes.append({"X": x, "Y": y, "Z": z})
            xs.append(x)
            ys.append(y)
            zs.append(z)

    with open(f"{f}.json", "w") as fh:
        json.dump(cubes, fh, indent=2)

    if not f.startswith("ex"):
        continue

    # build up the numpy logo
    n_voxels = np.zeros((max(xs) + 1, max(ys) + 1, max(zs) + 1), dtype=bool)
    for x, y, z in zip(xs, ys, zs):
        n_voxels[x, y, z] = True
    facecolors = np.where(n_voxels, "#FFD65DC0", "#7A88CCC0")
    edgecolors = np.where(n_voxels, "#BFAB6E", "#7D84A6")
    filled = np.ones(n_voxels.shape)

    # upscale the above voxel image, leaving gaps
    filled_2 = explode(n_voxels)
    fcolors_2 = explode(facecolors)
    ecolors_2 = explode(edgecolors)

    # Shrink the gaps
    x, y, z = np.indices(np.array(filled_2.shape) + 1).astype(float) // 2
    x[0::2, :, :] += 0.05
    y[:, 0::2, :] += 0.05
    z[:, :, 0::2] += 0.05
    x[1::2, :, :] += 0.95
    y[:, 1::2, :] += 0.95
    z[:, :, 1::2] += 0.95

    ax = plt.figure().add_subplot(projection="3d")
    ax.voxels(x, y, z, filled_2, facecolors=fcolors_2, edgecolors=ecolors_2)
    ax.set_aspect("equal")
    ax.set_xlabel("X")
    ax.set_ylabel("Y")
    ax.set_zlabel("Z")

    plt.show(block=True)
