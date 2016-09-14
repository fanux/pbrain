# -*- encoding: utf-8 -*-
import unittest

import requests
import json

HOST = "localhost"
PORT = "8081"

URL = "http://" + HOST + ":" + PORT 
session = requests.session()

"""
TODO 后期版本可能对插件优先级，策略优先级等进行定义
"""

class PluginTest(unittest.TestCase):
    def printJson(self, s):
        b = json.loads(s)
        print json.dumps(b, indent=4)

    def setUp(self):
        pass

    def tearDown(self):
        session.close()

    def test_example(self):
        data = {}
        response = session.post(URL + "", data = json.dumps(data))

        self.printJson(response.text)

    def test_createDecider(self):
        headers = {"content-type":"application/json"}
        data = {
            "Name":"plugin_decider",
            "Kind":"",
            "Status":"enable",
            "Description":"scaling container by time",
            "Spec":"""{"key":"value"}""",
            "Manual":"a markdown plugin document write manaul, how to config rules"
        }

        response = session.post(URL + "/plugins", headers=headers, data=json.dumps(data))


        # {"plugin":"time-plugin", "error":0}
        # print response.text
        self.printJson(response.text)

    def test_createPlugin(self):
        headers = {"content-type":"application/json"}
        data = {
            "Name":"plugin_pipeline",
            "Kind":"",
            "Status":"enable",
            "Description":"scaling container by time",
            "Spec":"""{"key":"value"}""",
            "Manual":"a markdown plugin document write manaul, how to config rules"
        }

        response = session.post(URL + "/plugins", headers=headers, data=json.dumps(data))


        # {"plugin":"time-plugin", "error":0}
        # print response.text
        self.printJson(response.text)

    def test_listPlugins(self):
        response = session.get(URL + "/plugins?offset=0&limit=20")

        """
        [
            {
                "Name":"time-plugin",
                "Kind":"",
                "Status":"enable",
                "Description":"scaling container by time",
                "Spec":"a json string"
            },
        ]
        """
        self.printJson(response.text)

    def test_getPluginDetail(self):
        plugin_name = "time-plugin"
        response = session.get(URL + "/plugins/" + plugin_name)

        """
         {
            "Name":"time-plugin",
            "Kind":"",
            "Status":"enable",
            "Description":"scaling container by time",
            "Spec":"a json string",
            "Manual":""
        }
        """
        self.printJson(response.text)

    def test_updatePlugin(self):
        plugin_name = "q-pipeline"
        headers = {"content-type":"application/json"}
        data =  {
            "Name":"q-pipeline",
            "Kind":"new newKind",
            "Status":"new new disable",
            "Description":"new new scaling container by time",
            "Spec":"new new a json string",
            "Manual":"new new manal new"
        }
        response = session.put(URL + "/plugins/" + plugin_name, headers=headers, data=json.dumps(data))

        """
        {
            "Name":"time-plugin",
            "Kind":"",
            "Status":"disable",
            "Description":"scaling container by time",
            "Spec":"a json string"
        }
        """
        self.printJson(response.text)

    def test_setPluginScope(self):
        plugin_name = "time-plugin"
        data = {
            "Name":"time-plugin",
            "NodeTags":{
                "Disk":"ssd",
                "Memery":"128G"
            },
            "HostNames":[
                "Dev-1-107",
                "Dev-1-108"
            ]
        }

        response = session.post("%s/plugins/%s/scope" % (URL, plugin_name), data=json.dumps(data))

        """
        """
        self.printJson(response.text)


    def test_getPluginScope(self):
        plugin_name = "time-plugin"

        response = session.get("%s/plugins/%s/scope" % (URL, plugin_name))

        """
        {
            "Name":"time-plugin",
            "NodeTags":{
                "Disk":"ssd",
                "Memery":"128G"
            },
            "HostNames":[
                "dev-1-107",
                "dev-1-108"
            ]
        }

        """
        self.printJson(response.text)

    def test_updatePluginScope(self):
        plugin_name = "time-plugin"
        data = {
            "Name":"time-plugin",
            "NodeTags":{
                "Disk":"ssd",
                "CPU":"64"
            },
            "HostNames":[
                "dev-1-109",
                "dev-1-110"
            ]
        }

        response = session.put("%s/plugins/%s/scope" % (URL, plugin_name), data=json.dumps(data))

        """
        {
            "Name":"time-plugin",
            "NodeTags":{
                "Disk":"ssd",
                "Memery":"128G"
            },
            "HostNames":[
                "dev-1-107",
                "dev-1-108"
            ]
        }

        """
        self.printJson(response.text)

    def test_createDeciderStrategy(self):
        plugin_name = "plugin_decider"

        self.document = [
    {
        "App":"ats",
        "Image":"ats:latest",
        "Priority":1,
        "MinNum":3,
        "Spec":[
            {
                "Metrical":[0, 20],
                "Number":-10,
            },
            {
                "Metrical":[20, 40],
                "Number":-5,
            },
            {
                "Metrical":[40, 60],
                "Number":0,
            },
            {
                "Metrical":[60, 80],
                "Number":5,
            },
            {
                "Metrical":[80, 100],
                "Number":10,
            }
        ]
    },
    {
        "App":"hadoop",
        "Image":"hadoop:latest",
        "Priority":2,
        "MinNum":2,
        "Spec":[
            {
                "Metrical":[0, 20],
                "Number":-10,
            },
            {
                "Metrical":[20, 40],
                "Number":-5,
            },
            {
                "Metrical":[40, 60],
                "Number":0,
            },
            {
                "Metrical":[60, 80],
                "Number":5,
            },
            {
                "Metrical":[80, 100],
                "Number":10,
            }
        ]
    },
    {
        "App":"redis",
        "Image":"redis:v1.0",
        "Priority":3,
        "MinNum":1,
        "Spec":[
            {
                "Metrical":[0, 20],
                "Number":-10,
            },
            {
                "Metrical":[20, 40],
                "Number":-5,
            },
            {
                "Metrical":[40, 60],
                "Number":0,
            },
            {
                "Metrical":[60, 80],
                "Number":5,
            },
            {
                "Metrical":[80, 100],
                "Number":10,
            }
        ]
    },
]
        plugin_strategy = {
            "PluginName":"plugin_decider",
            "Name":"ats-decider",
            "Status":"enable",
            #对controller来说就是一个字符串，不关心其内容
            "Document":json.dumps(self.document) 
        }

        headers = {"content-type":"application/json"}
        response = session.post("%s/plugins/%s/strategies" % (URL, plugin_name), 
                headers=headers,
                data=json.dumps(plugin_strategy))

        self.printJson(response.text)


    def test_createPluginStrategy(self):
        """
        不同的插件策略文档格式差异会非常大, 这里以时间为调度维度作个示例
        新的插件负责定义和解析策略文档并作出相应的动作

        xx_document的格式不定，所有惩罚列表都有一个 name 和 value

        同一个插件的策略文档格式相同，一个插件支持多个策略，仅是参数不同。
        Second Minute Hour Day Month Dayofweek 
          *      *     *    *   *      *
        时间配置遵循linux crontab
        """
        plugin_name = "time-plugin"

        self.document = [
            {    # 0点的时候启动20个ats 10个hadoop
                "Cron":"*/1 * * * * *",
                "Apps":[
                    {
                        "App":"ats",
                        "Number":20
                    },
                    {
                        "App":"hadoop:latest",
                        "Number":10
                    },
                ]
            }
        ]

        plugin_strategy = {
            "PluginName":"plugin_pipeline",
            "Name":"scale-by-hour",
            "Status":"enable",
            #对controller来说就是一个字符串，不关心其内容
            "Document":json.dumps(self.document) 
        }

        headers = {"content-type":"application/json"}
        response = session.post("%s/plugins/%s/strategies" % (URL, plugin_name), 
                headers=headers,
                data=json.dumps(plugin_strategy))

        """
        """
        self.printJson(response.text)

    def test_getPluginStrategies(self):
        plugin_name = "time-plugin"

        response = session.get("%s/plugins/%s/strategies" % (URL, plugin_name))

        """
        "Strategies":[
            {
                "PluginName":"time-plugin",
                "Name":"scale-by-hour",
                "Status":"enable",
                "Document":json.dumps(self.document) 
            },  
        ]
        """
        self.printJson(response.text)

    def test_getPluginStrategy(self):
        plugin_name = "time-plugin"
        strategy_name = "scale-by-hour"

        response = session.get("%s/plugins/%s/strategies/%s" % (URL, plugin_name, strategy_name))

        """
        {
            "PluginName":"time-plugin",
            "Name":"scale-by-hour",
            "Status":"enable",
            "Document":json.dumps(self.document) 
        }

        """
        self.printJson(response.text)

    def test_updatePluginStrategy(self):
        plugin_name = "time-plugin"
        strategy_name = "v"

        data = {
            "PluginName":"time-plugin",
            "Name":"v",
            "Status":"new enale",
            "Document":"new document",
        }

        headers = {"content-type":"application/json"}
        response = session.put("%s/plugins/%s/strategies/%s" % (URL, plugin_name, strategy_name),
                headers=headers,
                data=json.dumps(data))

        """
        """
        self.printJson(response.text)

    def test_deletePluginStrategy(self):
        plugin_name = "time-plugin"
        strategy_name = "scale-by-hour"

        response = session.delete("%s/plugins/%s/strategies/%s" % 
                (URL, plugin_name, strategy_name))

        """
        """
        self.printJson(response.text)

    def test_deletePlugin(self):
        plugin_name = "q-pipeline"

        response = session.delete(URL + "/plugins/" + plugin_name)

        """
        {
            "Msg": "delete ok", 
            "Code": "0"
        }
        """
        self.printJson(response.text)

    def test_scaleApps(self):
        data = [
            {
                "App":"172.20.1.128:5000/nginx:latest",
                "Number":2,
            },
            {
                "App":"172.20.1.128:5000/httpd:2.4",
                "Number":8,
            }
        ]

        headers = {"content-type":"application/json"}
        response = session.post("%s/plugins/scale" % URL,
                headers=headers,
                data=json.dumps(data))

        """
        """
        #self.printJson(response.text)

if __name__ == '__main__':
    unittest.main()
