CREATE EXTENSION aws_commons;
CREATE EXTENSION aws_lambda;

CREATE OR REPLACE FUNCTION invoke_lambda_function()
RETURNS void AS $$
BEGIN
  PERFORM aws_lambda.invoke_lambda('arn:aws:lambda:eu-west-1:491649323445:function:data-import-service', 
                                   'Payload', 
                                   'RequestResponse');
END;
$$ LANGUAGE plpgsql;
