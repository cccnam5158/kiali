import pytest
import tests.conftest as conftest

# Note: Number of services +1 Views Group Node
# Note: Node and Edge counts are based on traffic origainating from the Ingress
BOOKINFO_EXPECTED_NODES=8
BOOKINFO_EXPECTED_EDGES=7
BOOKINFO_EXPECTED_SERVICES = 4


PARAMS = {'graphType':'versionedApp', 'duration':'60s', 'injectServiceNodes':'false', 'edges':'trafficRatePerSecond'}

def test_service_graph_rest_endpoint(kiali_json):

    assert kiali_json != None, "Json: {}".format(kiali_json)

    # Validate that there are Nodes
    assert len(kiali_json.get('elements').get('nodes')) >= 1

    # Validate that there are Edges
    assert len(kiali_json.get('elements').get('edges')) >= 1

def test_service_graph_bookinfo_namespace_(kiali_client):
    PARAMS["namespaces"] = conftest.get_bookinfo_namespace()

    # Validate Node count
    nodes = kiali_client.graph_namespaces(params=PARAMS)["elements"]['nodes']
    #print ("Node count: {}".format(len(nodes)))
    assert len(nodes) >=  BOOKINFO_EXPECTED_NODES, "Expected Nodes: {}   Actual Nodes, {}".format(BOOKINFO_EXPECTED_NODES, len(nodes))

    # validate edge count
    edges = kiali_client.graph_namespaces(params=PARAMS)["elements"]['edges']
    #print ("Edge count: {}".format(len(edges)))
    assert len(edges) >= BOOKINFO_EXPECTED_EDGES, "Expected Edges: {}   Actual Edges, {}".format(BOOKINFO_EXPECTED_EDGES, len(nodes))
