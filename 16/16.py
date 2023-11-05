import json

ex = """
Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II
"""

in_ = """
Valve SW has flow rate=0; tunnels lead to valves LX, LD
Valve VS has flow rate=0; tunnels lead to valves JO, OO
Valve OO has flow rate=10; tunnels lead to valves KK, HD, VS, KI
Valve DZ has flow rate=8; tunnels lead to valves KV, GX, WQ, BA, PK
Valve GX has flow rate=0; tunnels lead to valves AA, DZ
Valve IF has flow rate=0; tunnels lead to valves OI, DW
Valve BO has flow rate=0; tunnels lead to valves UJ, ZT
Valve KI has flow rate=0; tunnels lead to valves OO, KU
Valve JT has flow rate=3; tunnels lead to valves FC, AM, KV, XP, XZ
Valve TQ has flow rate=0; tunnels lead to valves AA, DW
Valve KK has flow rate=0; tunnels lead to valves QW, OO
Valve NR has flow rate=0; tunnels lead to valves UG, XM
Valve VO has flow rate=0; tunnels lead to valves YR, AA
Valve MS has flow rate=17; tunnels lead to valves LT, LX
Valve JO has flow rate=0; tunnels lead to valves YR, VS
Valve ZB has flow rate=0; tunnels lead to valves UJ, LT
Valve ZT has flow rate=0; tunnels lead to valves XM, BO
Valve YR has flow rate=9; tunnels lead to valves VO, FY, WB, JO
Valve QS has flow rate=0; tunnels lead to valves QW, FY
Valve UD has flow rate=0; tunnels lead to valves CA, JB
Valve AP has flow rate=0; tunnels lead to valves CA, DW
Valve KV has flow rate=0; tunnels lead to valves JT, DZ
Valve JH has flow rate=0; tunnels lead to valves IK, UJ
Valve LD has flow rate=15; tunnels lead to valves IK, SW
Valve XK has flow rate=0; tunnels lead to valves XZ, BH
Valve XM has flow rate=11; tunnels lead to valves XP, CJ, ZT, NR
Valve FY has flow rate=0; tunnels lead to valves YR, QS
Valve GI has flow rate=22; tunnel leads to valve TI
Valve JB has flow rate=14; tunnels lead to valves WB, UD, WQ, HD
Valve DW has flow rate=6; tunnels lead to valves AP, TQ, NQ, IF, PK
Valve UJ has flow rate=13; tunnels lead to valves JH, ZB, BO
Valve KU has flow rate=0; tunnels lead to valves CA, KI
Valve WQ has flow rate=0; tunnels lead to valves JB, DZ
Valve BA has flow rate=0; tunnels lead to valves BH, DZ
Valve AA has flow rate=0; tunnels lead to valves YX, TQ, VO, GX, QP
Valve TI has flow rate=0; tunnels lead to valves GI, UG
Valve FC has flow rate=0; tunnels lead to valves QP, JT
Valve CA has flow rate=18; tunnels lead to valves KU, UD, AP
Valve QW has flow rate=25; tunnels lead to valves QS, KK
Valve XZ has flow rate=0; tunnels lead to valves JT, XK
Valve YX has flow rate=0; tunnels lead to valves AA, CJ
Valve OI has flow rate=0; tunnels lead to valves IF, BH
Valve NQ has flow rate=0; tunnels lead to valves AM, DW
Valve QP has flow rate=0; tunnels lead to valves AA, FC
Valve AM has flow rate=0; tunnels lead to valves NQ, JT
Valve XP has flow rate=0; tunnels lead to valves XM, JT
Valve BH has flow rate=12; tunnels lead to valves BA, XK, OI
Valve HD has flow rate=0; tunnels lead to valves OO, JB
Valve LT has flow rate=0; tunnels lead to valves MS, ZB
Valve LX has flow rate=0; tunnels lead to valves MS, SW
Valve CJ has flow rate=0; tunnels lead to valves XM, YX
Valve PK has flow rate=0; tunnels lead to valves DW, DZ
Valve IK has flow rate=0; tunnels lead to valves LD, JH
Valve WB has flow rate=0; tunnels lead to valves YR, JB
Valve UG has flow rate=21; tunnels lead to valves TI, NR
"""

for name, lines in zip(['ex', 'in'], [ex, in_]):
    res = []
    for line in lines.splitlines():
        if line:
            node = line[6:8]
            flow = int(line.split("=")[1].split(";")[0].strip())
            if "valves" in line:
                neighbors = line.split("valves")[1].split(",")
            else:
                neighbors = line.split("valve")[1].split(",")
            neighbors = [s.strip() for s in neighbors]
            print(node, flow, neighbors)
            res.append({"name": node, "flow": flow, "neighbors": neighbors})

    with open(f'{name}.json', 'w') as fh:
        json.dump(res, fh, indent=2)
