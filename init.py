#! /usr/bin/env python3

import os
import platform
import shutil
import sys
import warnings
from pylxd import Client

client = Client()

warnings.filterwarnings("ignore")

def convert(num):
    unit = 1000.0
    for x in ['', 'KB', 'MB', 'GB', 'TB']:
        if num < unit:
            return "%.0f%s" % (num, x)
        num /= unit

def get_driver():
    machine = platform.machine().lower()
    if machine.startswith(('arm', 'aarch64')):
        return "btrfs"
    return "zfs"

def create_storage():
    snap_data = os.environ['SNAP_DATA']
    free = ( shutil.disk_usage(snap_data)[-1] / 10 ) * 6 
    
    config = { "config": { "size": convert(free) },
            "driver": get_driver(), "name": "default" }

    try:
        client.storage_pools.get('default')
    except:
        try:
            client.storage_pools.create(config)
        except Exception as ex:
            print(ex)
            sys.exit(1)

def init_image(name):
    try:
        if client.images.get_by_alias(name):
          print('Image: ' + name + ' already exists')
          return
    except:
        print('Creating master image: ' + name)
        
    image = client.images.create_from_simplestreams(
        'https://cloud-images.ubuntu.com/daily',
        name)
    image.add_alias(name, '')

def main():
    create_storage()
    for img in [ 'bionic', 'xenial' ]:
        init_image(img)

main()
