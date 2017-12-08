#include <stdio.h>
#include <cstdio>
#include <stdlib.h>

#include <jack/jack.h>
#include <jack/types.h>
#include <jack/session.h>

jack_client_t *client;

int main()
{
	const char *client_name = "test";
	jack_status_t status;

	client = jack_client_open(client_name, JackSessionID, &status);
	if (client == NULL)
	{
		fprintf(stderr, "jack_client_open() failed, "
				"status = 0x%2.0x\n", status);

		if (status & JackServerFailed)
		{
			fprintf(stderr, "Unable to connect to JACK server\n");
		}

		exit(1);
	}

	fprintf(stderr, "jack_client_open() OK, "
			"status = 0x%2.0x\n", status);

	return 0;
}
