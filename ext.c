#include "sqlite3ext.h"
SQLITE_EXTENSION_INIT1

#include "goext.h"

__declspec(dllexport)
int sqlite3_go_init(
    sqlite3 *db,
    char **pzErrMsg,
    sqlite3_api_routines *pApi
) {
	SQLITE_EXTENSION_INIT2(pApi);

    int rc = SQLITE_OK;

    rc = sqlite3_create_function(db, "jaro", 2, SQLITE_UTF8, (void*)pApi, Jaro, 0, 0);
    rc = sqlite3_create_function(db, "regex", 2, SQLITE_UTF8, (void*)pApi, Regex, 0, 0);
    rc = sqlite3_create_function(db, "parsetime", 4, SQLITE_UTF8, (void*)pApi, ParseTime, 0, 0);

    return rc;
}