package io.fission;

import java.io.File;
import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.URLClassLoader;
import java.util.Enumeration;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;
import java.util.logging.Level;
import java.util.logging.Logger;

import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.Request;
import javax.ws.rs.Consumes;
import javax.ws.rs.core.MediaType;

import io.fission.Function;

public class Server {

  private Function fn;

  private static final int CLASS_LENGTH = 6;

  private static Logger logger = Logger.getGlobal();

  @Path("/")
  Response home(Request req) {
  	if (fn == null) {
      return Response
                .status(Response.Status.BAD_REQUEST)
                .entity( "Container not specialized")
                .type( MediaType.TEXT_PLAIN)
                .build();
  	} else {
      return ((Function) fn).call(req, null);
  	}
  }

  @POST
  @Path("/v2/specialize")
  @Consumes(MediaType.APPLICATION_JSON)
  Response specialize(FunctionLoadRequest flr) {
    long startTime = System.nanoTime();
    File file = new File(flr.getFilepath());
    if (!file.exists()) {
      return Response
                .status(Response.Status.BAD_REQUEST)
                .entity( "/userfunc/user not found")
                .type( MediaType.TEXT_PLAIN)
                .build();
    }

    String entryPoint = flr.getFunctionName();
    logger.log(Level.INFO, "Entrypoint class:" + entryPoint);
    if (entryPoint == null) {
      return Response
                .status(Response.Status.BAD_REQUEST)
                .entity("Entrypoint class is missing in the function")
                .type( MediaType.TEXT_PLAIN)
                .build();
    }

    JarFile jarFile = null;
    ClassLoader cl = null;
    try {

      jarFile = new JarFile(file);
      Enumeration<JarEntry> e = jarFile.entries();
      URL[] urls = { new URL("jar:file:" + file + "!/") };

      // TODO Check if the classloading can be improved for ex. use something like:
      // Thread.currentThread().setContextClassLoader(cl);
      if (this.getClass().getClassLoader() == null) {
        cl = URLClassLoader.newInstance(urls);
      } else {
        cl = URLClassLoader.newInstance(urls, this.getClass().getClassLoader());
      }

      if (cl == null) {
        return Response
                .status(Response.Status.BAD_REQUEST)
                .entity("Failed to initialize the classloader")
                .type( MediaType.TEXT_PLAIN)
                .build();
      }

      // Load all dependent classes from libraries etc.
      while (e.hasMoreElements()) {
        JarEntry je = e.nextElement();
        if (je.isDirectory() || !je.getName().endsWith(".class")) {
          continue;
        }
        String className = je.getName().substring(0, je.getName().length() - CLASS_LENGTH);
        className = className.replace('/', '.');
        cl.loadClass(className);
      }

      // Instantiate the function class
      fn = (Function) cl.loadClass(entryPoint).newInstance();

    } catch (MalformedURLException e) {
      e.printStackTrace();
      return Response
                .status(Response.Status.BAD_REQUEST)
                .entity("Error loading the Function class file")
                .type( MediaType.TEXT_PLAIN)
                .build();
    } catch (ClassNotFoundException e) {
      e.printStackTrace();
      return Response
                .status(Response.Status.BAD_REQUEST)
                .entity("Error loading Function or dependent class")
                .type( MediaType.TEXT_PLAIN)
                .build();
    } catch (InstantiationException e) {
      e.printStackTrace();
      return Response
                .status(Response.Status.BAD_REQUEST)
                .entity("Error creating a new instance of function class")
                .type( MediaType.TEXT_PLAIN)
                .build();
    } catch (IllegalAccessException e) {
      e.printStackTrace();
      return Response
                .status(Response.Status.BAD_REQUEST)
                .entity("Error creating a new instance of function class")
                .type( MediaType.TEXT_PLAIN)
                .build();
    } catch (IOException e) {
      e.printStackTrace();
      return Response
                .status(Response.Status.BAD_REQUEST)
                .entity("Error reading the JAR file")
                .type( MediaType.TEXT_PLAIN)
                .build();
    } finally {
      try {
        // cl.close();
        jarFile.close();
      } catch (IOException e) {
        e.printStackTrace();
        return Response
                .status(Response.Status.BAD_REQUEST)
                .entity("Error closing the file while loading the class")
                .type( MediaType.TEXT_PLAIN)
                .build();
      }
    }
    long elapsedTime = System.nanoTime() - startTime;
    logger.log(Level.INFO, "Specialize call done in: " + elapsedTime / 1000000 + " ms");
    return Response.ok("Done").build();
  }

}