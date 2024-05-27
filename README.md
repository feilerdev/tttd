# Todo to Tech Debt (ttd)

The aim of this tool is to help SWE Teams in the task of better controlling it's technical debts (td) by
automatizing the extraction of SATDs from the code to the used issues manager.

## How it works
    - Detect: searchs the code for a key-value pair, full-filled with td type and description, pre-accorded tag.
    - Transform: exports the key-value to json object.
    - Store: creates a json file with a array of objects.
    - Export: ?

