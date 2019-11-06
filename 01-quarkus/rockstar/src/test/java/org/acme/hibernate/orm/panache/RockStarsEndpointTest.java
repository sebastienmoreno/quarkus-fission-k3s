package org.acme.hibernate.orm.panache;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.containsString;
import static org.hamcrest.core.IsNot.not;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

@QuarkusTest
public class RockStarsEndpointTest {

    @Test
    public void testListAllRockStars() {
        //List all, should have all 3 fruits the database has initially:
        given()
              .when().get("/")
              .then()
              .statusCode(200)
              .body(
                    containsString("Prince"),
                    containsString("Magma"),
                    containsString("Bruce Springsteen")
                    );
    }

}
