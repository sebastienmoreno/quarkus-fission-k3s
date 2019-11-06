package org.acme.hibernate.orm.panache;

import javax.persistence.Cacheable;
import javax.persistence.Column;
import javax.persistence.Entity;

import io.quarkus.hibernate.orm.panache.PanacheEntity;

@Entity
@Cacheable
public class RockStar extends PanacheEntity {

    @Column(length = 40, unique = true)
    public String name;

    public RockStar() {
    }

    public RockStar(String name) {
        this.name = name;
    }
}
