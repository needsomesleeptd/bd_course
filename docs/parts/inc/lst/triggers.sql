
CREATE OR REPLACE FUNCTION addToQueue()
Returns TRIGGER
AS $$
BEGIN
    IF NEW.document_name NOT LIKE '%VIP%' 
    THEN
        INSERT INTO 
        document_queues(doc_id,status)
        VALUES
        (NEW.ID,0); -- status zero means unchecked
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insertQueueTrigger AFTER INSERT ON public.documents
FOR EACH ROW EXECUTE PROCEDURE addToQueue();