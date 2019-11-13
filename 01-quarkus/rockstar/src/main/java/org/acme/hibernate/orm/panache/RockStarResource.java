package org.acme.hibernate.orm.panache;

import java.util.List;

import javax.enterprise.context.ApplicationScoped;
import javax.json.Json;
import javax.transaction.Transactional;
import javax.ws.rs.Consumes;
import javax.ws.rs.DELETE;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.WebApplicationException;
import javax.ws.rs.core.Response;
import javax.ws.rs.ext.ExceptionMapper;
import javax.ws.rs.ext.Provider;

import org.jboss.resteasy.annotations.jaxrs.PathParam;

import io.quarkus.panache.common.Sort;

@Path("/")
@ApplicationScoped
@Produces("application/json")
@Consumes("application/json")
public class RockStarResource {

    @GET
    public List<RockStar> get() {
        return RockStar.listAll(Sort.by("name"));
    }

    @POST
    @Transactional
    public Response create(RockStar RockStar) {
        if (RockStar.id != null) {
            throw new WebApplicationException("Id was invalidly set on request.", 422);
        }

        RockStar.persist();
        return Response.ok(RockStar).status(201).build();
    }

    @PUT
    @Path("{id}")
    @Transactional
    public RockStar update(@PathParam Long id, RockStar RockStar) {
        if (RockStar.name == null) {
            throw new WebApplicationException("RockStar Name was not set on request.", 422);
        }

        RockStar entity = RockStar.findById(id);

        if (entity == null) {
            throw new WebApplicationException("RockStar with id of " + id + " does not exist.", 404);
        }

        entity.name = RockStar.name;

        return entity;
    }

    @DELETE
    @Path("{id}")
    @Transactional
    public Response delete(@PathParam Long id) {
        RockStar entity = RockStar.findById(id);
        if (entity == null) {
            throw new WebApplicationException("RockStar with id of " + id + " does not exist.", 404);
        }
        entity.delete();
        return Response.status(204).build();
    }

    @Provider
    public static class ErrorMapper implements ExceptionMapper<Exception> {

        @Override
        public Response toResponse(Exception exception) {
            int code = 500;
            if (exception instanceof WebApplicationException) {
                code = ((WebApplicationException) exception).getResponse().getStatus();
            }
            return Response.status(code)
                    .entity(Json.createObjectBuilder().add("error", exception.getMessage()).add("code", code).build())
                    .build();
        }

    }
}
