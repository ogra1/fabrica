#! /usr/bin/env python3

import os
import sys
import threading
import time
import warnings

from pylxd import Client
client = Client()

warnings.filterwarnings("ignore")

thread_errors = []

def cleanup(container, exitcode, exception):
    print('Stopping container...')
    container.stop(wait=True)
	
    print('Deleting container')
    container.delete()
	
    if exception:	
        print(exception)
    sys.exit(exitcode)

def parse_args(tree):
    dir = tree.rsplit('/')[-1]

    if tree.endswith('.git'):
        dir = tree.rsplit('/')[-1][:-len('.git')]

    cname = dir + '-' + str(int(time.time()))

    return cname, dir

def launch_container(cname, release):
    config = {'name': cname, 'source': {'type': 'image', 
        'alias': release}}

    container = client.containers.create(config, wait=True)

    print('Starting container ' + cname)
    container.start(wait=True)

    print('Waiting for network...')
    for i in range(1, 15):
        if container.execute(['ping', '-c1', '8.8.8.8']) == 0:
            break

    return container			

def run_in_container(cmds, container, outbuf, errbuf):
    for command in cmds:
        try:
            if not len(thread_errors):
                print('\nRunning command: ' + ' '.join(command))
                out = container.execute(command,
                    environment = { 'FLASH_KERNEL_SKIP': 'true',
                        'DEBIAN_FRONTEND': 'noninteractive',
                        'TERM': 'xterm',
                        'SNAPCRAFT_BUILD_ENVIRONMENT': 'host'},
                        stdout_handler=outbuf.append,
                        stderr_handler=errbuf.append)
        except:
            pass

        if out.exit_code != 0:
            thread_errors.append('Command failed: ' + ' '.join(command))		

def spawn_async(cmd, container):
    stdout_buffer = []
    stderr_buffer = []

    thread = threading.Thread(target=run_in_container,
        args=[cmd, container, stdout_buffer, stderr_buffer])

    thread.start()

    while thread.isAlive():
        if len(thread_errors) > 0:
            thread.do_run = False
            raise Exception(str(thread_errors[0]))	
        while len(stdout_buffer):
            print(stdout_buffer.pop(0).strip())		
        while len(stderr_buffer):
            print(stderr_buffer.pop(0).strip())		
    
    thread.join()

def pull_snap(dir, container, build_id):
    snap = container.execute(['sh', '-c', 'ls /root/' + dir + '/*.snap']).stdout.rstrip()
    out_dir = os.environ["HOME"] + '/' + build_id
    os.makedirs(out_dir, exist_ok=True)
    outsnap = out_dir + '/' + os.path.basename(snap)
    print('Pulling snap package ' + snap + ' to ' + outsnap)
    try:
        filedata = container.files.get(snap)
        filewr = open(outsnap, 'wb')
        filewr.write(filedata)
        filewr.close()
        print('Archived snap package: ' + outsnap)
    except Exception as ex:
        print('can not write ' + outsnap)
        print(ex)

def main():
    try:
        tree = sys.argv[1]
        build_id = sys.argv[2]
        distro = sys.argv[3]
        cname, dir = parse_args(tree)

        container = launch_container(cname, distro)

        prepare = [['apt', 'update'],
                   ['apt', '-y', 'upgrade'],
                   ['apt', '-y', 'install', 'build-essential'],
                   ['apt', '-y', 'clean'],
                   ['snap', 'install', 'snapcraft', '--classic']]

        print('Preparing container for build...')
        spawn_async(prepare, container)		

        build = [['git', 'clone', '--progress', tree],
                 ['sh', '-cv', 'cd /root/' + dir + '; snapcraft']]

        print('Running build for ' + tree)
        spawn_async(build, container)		

        pull_snap(dir, container, build_id)
    except KeyboardInterrupt:
            cleanup(container, 0, '')
            raise
    except Exception as ex:
            cleanup(container, 1, ex)
    cleanup(container, 0, '')

main()
