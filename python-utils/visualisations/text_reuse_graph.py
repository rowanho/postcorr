import json
import sys

import matplotlib.pyplot as plt
import networkx as nx

def build_graph(filepath):
    G = nx.Graph()
    with open(filepath) as json_file:
        data = json.load(json_file)
        for docId in data:
            for e in data[docId]:
                G.add_edge(docId, e['docId'], weight=e['score'])
    
    return G
    
if __name__ == "__main__":
    filepath = sys.argv[1]
    G = build_graph(filepath)
    pos = nx.fruchterman_reingold_layout(G)
    nx.draw(G, pos)
    #labels = nx.get_edge_attributes(G, 'weight')
    #nx.draw_networkx_edge_labels(G, pos, labels=labels)
    plt.show()