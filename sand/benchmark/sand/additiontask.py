import pika
import json

# @nuclio.configure
#
# function.yaml:
#   apiVersion: "nuclio.io/v1"
#   kind: NuclioFunction
#   metadata:
#     name: additiontask
#     namespace: nuclio
#   spec:
#     env:
#     - name: SAND_PORT_INTERNAL
#       value: "10000"
#     handler: main:handler
#     runtime: python
#     triggers:
#       http:
#         kind: http
#         attributes:
#           serviceType: ClusterIP

def handler(context, event):

    broker_url = "amqp://jeff:jeff@localhost:5672/%2F"
    
    def invoke(function_name: str, payload: any) -> None:
        connection = pika.BlockingConnection(pika.URLParameters(broker_url))
        channel = connection.channel()
        channel.basic_publish(exchange="", routing_key=function_name, body=payload)
        connection.close()
  
    payload = json.loads(event.body)
    target, actual = int(payload["target"]), int(payload["actual"])

    if actual >= target:
        invoke("entrytask_result", f"{actual}")

    else:
        payload = json.dumps({"target": target, "actual": actual+1})
        invoke("additiontask", payload)

    return context.Response(
        body="",
        content_type='text/plain',
        status_code=200
    )
