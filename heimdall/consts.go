package heimdall

var queryCreateOrReplacePublishFunction = `
	CREATE OR REPLACE FUNCTION %v() RETURNS TRIGGER AS $$
    DECLARE 
        data json;
        notification json;
    
    BEGIN
		-- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;
        
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        
                        
        PERFORM pg_notify('%v', notification::text);
        RETURN NULL; 
    END;
    
	$$ LANGUAGE plpgsql;`

var queryCreateTriggerStatement = `
	CREATE TRIGGER %s
	AFTER INSERT OR UPDATE OR DELETE ON %s
	FOR EACH ROW EXECUTE PROCEDURE %s();
`

var queryDropTriggerStatement = `
	DROP TRIGGER IF EXISTS %s ON %s;
`
