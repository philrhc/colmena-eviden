/*
  COLMENA-DESCRIPTION-SERVICE
  Copyright © 2024 EVIDEN
  
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  
  http://www.apache.org/licenses/LICENSE-2.0
  
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  
  This work has been implemented within the context of COLMENA project.
*/
/*
Copyright © 2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This work has been implemented within the context of COLMENA project.
*/
package main

import (
	"context"
	"log"
	"time"

	pb "colmena/sla-management-svc/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	addr := "localhost:8099"

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSLAsvcClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slao := &pb.InputSLAObj{
		ServiceId: "test_grpc_02",
		Kpis:      []string{},
		Roles: []*pb.InputSLAObj_InputSLAObjRole{
			{
				Id: "test_grpc_02_Plantwatering",
				Kpis: []string{
					"go_memstats_frees_total < 50000",
				},
			},
		},
	}

	r, err := c.CreateSLA(ctx, slao)
	if err != nil {
		log.Fatalf("could not create SLA: %v", err)
	}
	log.Println(r)

}
