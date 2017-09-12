#!/usr/bin/env python3
# generate
# Script to generate test data for the stats package.
#
# Author:  Benjamin Bengfort <benjamin@bengfort.com>
# Created: Tue Sep 12 09:35:41 2017 -0400
#
# ID: generate.py [] benjamin@bengfort.com $

"""
Script to generate test data for the stats package.
"""

##########################################################################
## Imports
##########################################################################

import os
import json
import argparse
import numpy as np


BASEDIR = os.path.dirname(__file__)
STANDARDIZED = os.path.join(BASEDIR, "standardized.txt")
LATENCIES = os.path.join(BASEDIR, "latencies.txt")


##########################################################################
## Helper Classes
##########################################################################

class NumpyEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, np.integer):
            return int(obj)
        elif isinstance(obj, np.floating):
            return float(obj)
        elif isinstance(obj, np.ndarray):
            return obj.tolist()
        else:
            return super(NumpyEncoder, self).default(obj)


def load(path):
    if not os.path.exists(path):
        return

    with open(path, 'r') as f:
        for line in f:
            line = line.strip()
            line = line.strip("ms")
            yield np.float64(line)


##########################################################################
## Commands
##########################################################################

def generate_standardized(n=1000000):
    if os.path.exists(STANDARDIZED):
        print("standardized dataset already exists ... skipping")

    with open(STANDARDIZED, "w") as f:
        for s in np.random.normal(0, 1, n):
            f.write("{}ms\n".format(s))

def generate_latencies(n=1000000):
    if os.path.exists(STANDARDIZED):
        print("standardized dataset already exists ... skipping")
        
    with open(LATENCIES, "w") as f:
        for s in np.random.normal(120, 17, n):
            f.write("{}ms\n".format(s))

def report_standardized():
    data = np.array(list(load(STANDARDIZED)))
    if len(data) < 1: return

    stats = {
        "minimum": data.min(),
        "maximum": data.max(),
        "range": data.max() - data.min(),
        "samples" : data.shape[0],
        "mean": data.mean(),
        "stddev": data.std(),
        "variance": data.var(),
    }

    print("Standardized:")
    print(json.dumps(stats, indent=2, cls=NumpyEncoder))
    print

def report_latencies():
    data = np.array(list(load(LATENCIES)))
    if len(data) < 1: return
    data = data * 1000000
    data = data.astype(np.int64)

    stats = {
        "fastest": data.min(),
        "slowest": data.max(),
        "range": data.max() - data.min(),
        "samples" : data.shape[0],
        "mean": data.mean(),
        "stddev": data.std(),
        "variance": data.var(),
        "throughput": data.shape[0] / (data.sum() * 1e-9),
    }

    print("Latencies:")
    print(json.dumps(stats, indent=2, cls=NumpyEncoder))


##########################################################################
## Main Method
##########################################################################

if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description="generate test data for the stats package",
    )

    parser.add_argument(
        "-s", "--standardized", action="store_true",
        help="generate standardized data"
    )

    parser.add_argument(
        "-l", "--latencies", action="store_true",
        help="generate latencies data"
    )

    parser.add_argument(
        "-n", "--samples", type=int, metavar="N", default=1000000,
        help="number of samples to generate"
    )

    # Parse the arguments from the command line
    args = parser.parse_args()

    # Generate the standardized data set if requested, and report the values
    # if the standardized dataset exists.
    if args.standardized:
        generate_standardized(args.samples)
    report_standardized()

    # Generate the latencies data set if requested, and report the values
    # if the latencies dataset exists.
    if args.latencies:
        generate_latencies(args.samples)
    report_latencies()
