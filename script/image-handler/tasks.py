import json
from invoke import task,run

"""
* default configure file is `config`
* `$invoke [function name] --config [configure file]`
* `invoke add_node_tags --config host_info.json`
* `invoke clean_node_tags --config host_info.json`
"""

def get_config(config):
    f = open(config)
    config = json.load(f)

    f.close()
    return config

@task
def save(ctx, config="config"):
    config = get_config(config)

    images = config["images"]

    for i in images:
	ctx.run("docker save %s > images/%s.tar" % (i,i))

@task
def load(ctx, config="config"):
    config = get_config(config)

    images = config["images"]

    for i in images:
	ctx.run("docker load -i images/%s.tar" % i)
