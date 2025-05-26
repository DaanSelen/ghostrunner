#!/bin/python3
#
#
from configparser import ConfigParser
import os

#
#
#

class utilities:
    @staticmethod
    def load_config(segment: str = 'ghostrunner') -> dict:
        '''
        Function that loads the segment from the config.conf (by default) file and returns the it in a dict.
        '''

        conf_file = "./conf/ghostserver.conf"
        if not os.path.exists(conf_file):
            print(f'Missing config file {conf_file}. Provide an alternative path.')
            os._exit(1)

        config = ConfigParser()
        try:
            config.read(conf_file)
        except Exception as err:
            print(f"Error reading configuration file '{conf_file}': {err}")
            os._exit(1)

        if segment not in config:
            print(f'Segment "{segment}" not found in config file {conf_file}.')
            os._exit(1)

        return dict(config[segment])
