import pika
import json

# @nuclio.configure
#
# function.yaml:
#   apiVersion: "nuclio.io/v1"
#   kind: NuclioFunction
#   metadata:
#     name: entrytask
#     namespace: nuclio
#   spec:
#     env:
#     - name: SAND_PORT_INTERNAL
#       value: "8080"
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
        target, actual = int(event.body), 0
    
        payload = json.dumps({"target": target, "actual": actual})
        invoke("additiontask", payload)

    def wait_for_result():

        connection = pika.BlockingConnection(pika.URLParameters(broker_url))
        channel = connection.channel()
        channel.queue_declare(queue="entrytask_result")

        done = False
        msg = None

        def callback(channel, method, properties, body):
            # `callback takes one position argument but four were given` => hat irgendwie 4 params auch wenn ich 3 nicht benutze >.<
            nonlocal done, msg
            msg = body
            done = True

        channel.basic_consume(queue="entrytask_result", on_message_callback=callback, auto_ack=True)

        while not done:
            connection.process_data_events()

        connection.close()
        return msg

    return context.Response(
        body=wait_for_result(),
        content_type='text/plain',
        status_code=200
    )
