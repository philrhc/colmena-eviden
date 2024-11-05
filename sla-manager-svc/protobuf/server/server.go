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
package server

import (
	context "context"
	"fmt"
	"net"
	"strconv"
	"sync"

	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"
	pb "colmena/sla-management-svc/protobuf"

	uuid "github.com/lithammer/shortuuid/v4"
	"google.golang.org/grpc"
)

// path used in logs
const pathLOG string = "SLA > gRPC > "

var repository model.IRepository

type server_sla struct {
	pb.UnimplementedSLAsvcServer
}

/**
 * CreateSLA creates a new SLA
 */
func (s *server_sla) CreateSLA(ctx context.Context, in *pb.InputSLAObj) (*pb.InputSLAObj, error) {
	logs.GetLogger().Info(pathLOG + "<< CreateSLA >>")
	logs.GetLogger().Debug(pathLOG + "CreateSLA > SLAObj: " + in.String())

	slas, _ := slaObtToSlaModel(in)

	for _, sla := range slas {
		_, e := repository.CreateSLA(&sla)
		if e != nil {
			logs.GetLogger().Error(pathLOG + "Error creating SLA: " + e.Error())
		}
	}

	return in, nil
}

/**
 * DeleteSLA deletes a SLA
 */
func (s *server_sla) DeleteSLA(ctx context.Context, in *pb.SLAId) (*pb.SLAObj, error) {
	logs.GetLogger().Info(pathLOG + "<< DeleteSLA >>")
	logs.GetLogger().Debug(pathLOG + "DeleteSLA > Id: " + in.Id)

	e := repository.DeleteSLA(in.Id)

	if e != nil {
		return nil, e
	}

	return nil, nil
}

func (s *server_sla) GetSLA(ctx context.Context, in *pb.SLAId) (*pb.SLAObj, error) {
	logs.GetLogger().Info(pathLOG + "<< GetSLA >>")
	logs.GetLogger().Debug(pathLOG + "GetSLA > Id: " + in.Id)

	m, e := repository.GetSLA(in.Id)

	qos, _ := slaModelToSlaObj(*m) //(pb.SLAObj, error)

	if e != nil {
		return nil, e
	}

	return &qos, nil
}

/*
InitializegRPCServer initialization function
*/
func InitializegRPCServer(wg *sync.WaitGroup, r model.IRepository) {
	logs.GetLogger().Info(pathLOG + "[InitializegRPCServer] Initializing gRPC Server ...")

	CreateServer(wg, "8099", r)
}

/**
 * CreateServer creates a gRPC server
 */
func CreateServer(wg *sync.WaitGroup, port string, r model.IRepository) {

	repository = r

	logs.GetLogger().Info(pathLOG + "[Create Server] Creating gRPC server ...")

	defer wg.Done()

	p, err := strconv.Atoi(port)

	if err != nil {
		logs.GetLogger().Fatal(pathLOG+"[Create Server] Invalid port: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", p))

	if err != nil {
		logs.GetLogger().Fatal(pathLOG+"[Create Server] failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterSLAsvcServer(s, &server_sla{})

	logs.GetLogger().Infof(pathLOG+"[Create Server] gRPC SLA server listening at ", lis.Addr())

	if err := s.Serve(lis); err != nil {
		logs.GetLogger().Fatal(pathLOG+"[Create Server] failed to serve: %v", err)
	}

}

// slaObtToSlaModel
func slaObtToSlaModel(input *pb.InputSLAObj) ([]model.SLA, error) {

	var slas []model.SLA

	// InputSLA ==> SLA(s)
	if len(input.Roles) > 0 {
		for _, r := range input.Roles {
			if len(r.Kpis) > 0 {
				uid := uuid.New()
				sla := model.SLA{}

				sla.Name = input.ServiceId
				sla.Id = input.ServiceId + "-" + uid
				sla.State = "started"

				sla.Details.Guarantees = make([]model.Guarantee, 1) // TODO for each KPI => 1 Guarantee
				sla.Details.Guarantees[0].Name = r.Id
				sla.Details.Guarantees[0].Constraint = r.Kpis[0]

				slas = append(slas, sla)
			}
		}
	}

	return slas, nil

}

// slaModelToSlaObj
func slaModelToSlaObj(in model.SLA) (pb.SLAObj, error) {

	sla := pb.SLAObj{}

	sla.Name = in.Name
	sla.Id = in.Id
	sla.State = "started"

	if len(in.Details.Guarantees) > 0 {
		sla.Details = make([]*pb.SLAObj_Details, 1)

		sla.Details[0].Guarantees[0].Name = in.Details.Guarantees[0].Name
		sla.Details[0].Guarantees[0].Constraint = in.Details.Guarantees[0].Constraint

	}

	return sla, nil

}
